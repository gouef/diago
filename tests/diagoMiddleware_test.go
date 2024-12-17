package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/gouef/diago"
	"github.com/stretchr/testify/assert"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MyExtension struct{}

func (e *MyExtension) GetPanelHtml(c *gin.Context) string {
	return "<div>Panel</div>"
}

func (e *MyExtension) GetHtml(c *gin.Context) string {
	return "<div>Content</div>"
}

func (e *MyExtension) GetJSHtml(c *gin.Context) string {
	return "<script>console.log('JS');</script>"
}

func (e *MyExtension) BeforeNext(c *gin.Context) {
}

func (e *MyExtension) AfterNext(c *gin.Context) {
}

func TestDiagoMiddleware(t *testing.T) {
	r := gin.Default()
	gin.SetMode(gin.TestMode)
	d := diago.NewDiago()

	d.AddExtension(&MyExtension{})

	r.GET("/test", diago.DiagoMiddleware(nil, d), func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, "Hello, world!")
	})

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))

	assert.Contains(t, w.Body.String(), "<div>Panel</div>")
	assert.Contains(t, w.Body.String(), "<div>Content</div>")
	assert.Contains(t, w.Body.String(), "<script>console.log('JS');</script>")
}

func TestGenerateDiagoPanelHTML(t *testing.T) {
	diagoData := diago.DiagoData{
		ExtensionsHtml:      []template.HTML{"<div>Extension HTML</div>"},
		ExtensionsPanelHtml: []template.HTML{"<div>Panel HTML</div>"},
		ExtensionsJSHtml:    []template.HTML{"<script>JS</script>"},
	}

	diagoPanelHTML, err := diago.GenerateDiagoPanelHTML(diagoData)

	assert.NoError(t, err)

	assert.Contains(t, diagoPanelHTML, "<div>Extension HTML</div>")
	assert.Contains(t, diagoPanelHTML, "<div>Panel HTML</div>")
	assert.Contains(t, diagoPanelHTML, "<script>JS</script>")
}
