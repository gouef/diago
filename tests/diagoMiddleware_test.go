package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/gouef/diago"
	"github.com/gouef/diago/extensions"
	"github.com/gouef/router"
	"github.com/stretchr/testify/assert"
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

func (e *MyExtension) BeforeNext(c *gin.Context) {}

func (e *MyExtension) AfterNext(c *gin.Context) {}

func (e *MyExtension) SetTemplateProvider(provider diago.TemplateProvider) {

}
func (e *MyExtension) GetTemplateProvider() diago.TemplateProvider {
	return extensions.NewDefaultTemplateProvider()
}
func (e *MyExtension) SetPanelGenerator(generator diago.PanelGenerator) {

}
func (e *MyExtension) GetPanelGenerator() diago.PanelGenerator {
	return diago.NewDefaultPanelGenerator()
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
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	r.ServeHTTP(w, req)

	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))

	assert.Contains(t, w.Body.String(), "<div>Panel</div>")
	assert.Contains(t, w.Body.String(), "<div>Content</div>")
	assert.Contains(t, w.Body.String(), "<script>console.log('JS');</script>")
}

func TestDiagoMiddlewareReleaseMode(t *testing.T) {
	r := router.NewRouter()
	r.EnableRelease()
	n := r.GetNativeRouter()
	gin.SetMode(gin.ReleaseMode)

	d := diago.NewDiago()

	d.AddExtension(&MyExtension{})

	n.GET("/test", diago.DiagoMiddleware(r, d), func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, "Hello, world!")
	})
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	n.ServeHTTP(w, req)

	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))

	assert.Equal(t, w.Body.String(), "Hello, world!")
}
