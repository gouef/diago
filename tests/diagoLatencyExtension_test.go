package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/gouef/diago/extensions"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDiagoLatencyExtension(t *testing.T) {
	latencyExtension := extensions.NewDiagoLatencyExtension()

	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {

		latencyExtension.BeforeNext(c)

		time.Sleep(100 * time.Millisecond)

		latencyExtension.AfterNext(c)

		panelHtml := latencyExtension.GetPanelHtml(c)
		c.String(http.StatusOK, panelHtml)
	})

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Contains(t, w.Body.String(), "ms")
}
