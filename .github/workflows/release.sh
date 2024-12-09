PREV_TAG=$(git describe --tags --abbrev=0 HEAD^ 2>/dev/null || echo "")
          if [ -z "$PREV_TAG" ]; then
            echo "No previous tag found, using initial commit."
            COMMITS=$(git log --format="- %s (@%an)" HEAD)
          else
            COMMITS="<ul>"
            while read -r commit_hash; do
              COMMIT_MSG=$(git log -n 1 --pretty=format:"%s" "$commit_hash")
              AUTHOR=$(gh api "/repos/${{ github.repository }}/commits/$commit_hash" --jq '.author.login')
              COMMITS+="<li>"
              # Pokud commit je fixup, připojíme ho k původnímu commit message
              if [[ "$COMMIT_MSG" == fixup!* ]]; then
                ORIGINAL_COMMIT_MSG=$(git log -n 1 --pretty=format:"%s" "$commit_hash"^)
                COMMITS+="$ORIGINAL_COMMIT_MSG (fixup) (@$AUTHOR)"
              else
                COMMITS+="$COMMIT_MSG (@$AUTHOR)"
              fi

              # Prohledáme všechny komentáře v PR/Issue pro zmínky na commit
              COMMENT_LINKS=""
              # Hledání v issue
              ISSUE_COMMENTS=$(gh api "/repos/${{ github.repository }}/issues/comments?per_page=100" --jq "[.[] | select(.body | contains(\"$commit_hash\"))]")
              if [[ -n "$ISSUE_COMMENTS" ]]; then
                # Pro každé nalezené komentáře přidáme odkaz na issue
                ISSUE_URLS=$(gh api "/repos/${{ github.repository }}/issues?state=all" --jq "[.[] | select(.body | contains(\"$commit_hash\"))][].html_url")
                for ISSUE_URL in $ISSUE_URLS; do
                  COMMENT_LINKS+="[Link]($ISSUE_URL), "
                done
              fi

              PR_LINKS=""
              # Hledání v pull requestech
              PR_COMMENTS=$(gh api "/repos/${{ github.repository }}/pulls/comments?per_page=100" --jq "[.[] | select(.body | contains(\"$commit_hash\"))]")
              if [[ -n "$PR_COMMENTS" ]]; then
                # Pro každé nalezené komentáře přidáme odkaz na PR
                PR_URLS=$(gh api "/repos/${{ github.repository }}/pulls?state=all" --jq "[.[] | select(.body | contains(\"$commit_hash\"))][].html_url")
                for PR_URL in $PR_URLS; do
                  PR_LINKS+="[Link]($PR_URL), "
                done
              fi

              if [[ -n "$PR_LINKS" ]]; then
                COMMITS+=" (Referenced in PRs: $PR_LINKS)"
              fi

              if [[ -n "$COMMENT_LINKS" ]]; then
                COMMITS+=" (Referenced in Issues: $COMMENT_LINKS)"
              fi
              COMMITS+="</li>"
            done < <(git log --format="%H" $PREV_TAG..HEAD)
            COMMITS+="</ul>"
          fi

          COMMITS=$(echo "$COMMITS" | sed 's/[[:cntrl:]]//g')

          echo "commits=$COMMITS" >> $GITHUB_ENV