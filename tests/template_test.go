package tests

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gouef/diago"
	"github.com/gouef/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"html/template"
	_ "log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockTemplateProvider struct {
	mock.Mock
}

func (m *MockTemplateProvider) GetTemplate() string {
	args := m.Called()
	return args.String(0)
}

type MockPanelGenerator struct {
}

func (m *MockPanelGenerator) GenerateHTML(name string, templateProvider diago.TemplateProvider, data interface{}) (string, error) {
	return "", errors.New("template parsing error")
}

func TestDiagoMiddleware_GenerateHTML_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockPanelGenerator := new(MockPanelGenerator)

	r := router.NewRouter()
	n := r.GetNativeRouter()
	n.LoadHTMLGlob("templates/*")

	middleware := diago.DiagoMiddleware(r, &diago.Diago{PanelGenerator: mockPanelGenerator, TemplateProvider: diago.NewDefaultTemplateProvider()})

	n.Use(middleware)

	r.AddRouteGet("notfound", "/notfound", func(c *gin.Context) {
		c.HTML(http.StatusOK, "status.gohtml", gin.H{"content": template.HTML("<div>OK</div>")})
	})

	t.Run("Test Custom 404 Handler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
		assert.Equal(t, `<div id="content"><div>OK</div></div>Error generating Diago panel HTML`, w.Body.String())
	})
}

func TestGenerateHTML_ParseError(t *testing.T) {
	mockTemplateProvider := new(MockTemplateProvider)
	mockTemplateProvider.On("GetTemplate").Return("{{.Invalid") // Nesprávná šablona, která způsobí chybu

	d := &diago.DefaultPanelGenerator{}

	result, err := d.GenerateHTML("diagoPanel", mockTemplateProvider, nil)

	assert.Empty(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unclosed action")
}