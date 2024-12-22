package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/gouef/diago"
	"github.com/gouef/diago/extensions"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockDiagoExtension struct {
	PanelHtml        string
	Html             string
	JSHtml           string
	BeforeCalled     bool
	AfterCalled      bool
	PanelGenerator   diago.PanelGenerator
	TemplateProvider diago.TemplateProvider
}

func (m *MockDiagoExtension) GetPanelHtml(c *gin.Context) string {
	return m.PanelHtml
}

func (m *MockDiagoExtension) GetHtml(c *gin.Context) string {
	return m.Html
}

func (m *MockDiagoExtension) GetJSHtml(c *gin.Context) string {
	return m.JSHtml
}

func (m *MockDiagoExtension) BeforeNext(c *gin.Context) {
	m.BeforeCalled = true
}

func (m *MockDiagoExtension) AfterNext(c *gin.Context) {
	m.AfterCalled = true
}

func (e *MockDiagoExtension) SetTemplateProvider(provider diago.TemplateProvider) {

}
func (e *MockDiagoExtension) GetTemplateProvider() diago.TemplateProvider {
	return extensions.NewDefaultTemplateProvider()
}
func (e *MockDiagoExtension) SetPanelGenerator(generator diago.PanelGenerator) {

}
func (e *MockDiagoExtension) GetPanelGenerator() diago.PanelGenerator {
	return diago.NewDefaultPanelGenerator()
}

func TestDiago(t *testing.T) {
	t.Run("Test AddExtension", func(t *testing.T) {
		newDiago := diago.NewDiago()

		mockExtension := &MockDiagoExtension{
			PanelHtml: "<div>Panel</div>",
		}

		newDiago.AddExtension(mockExtension)

		assert.Len(t, newDiago.GetExtensions(), 1)
		assert.Equal(t, "<div>Panel</div>", newDiago.GetExtensions()[0].GetPanelHtml(nil))
	})

	t.Run("Test GetExtensions", func(t *testing.T) {
		newDiago := diago.NewDiago()

		mockExtension1 := &MockDiagoExtension{PanelHtml: "<div>Panel1</div>"}
		mockExtension2 := &MockDiagoExtension{PanelHtml: "<div>Panel2</div>"}

		newDiago.AddExtension(mockExtension1).AddExtension(mockExtension2)

		assert.Len(t, newDiago.GetExtensions(), 2)
		assert.Equal(t, "<div>Panel1</div>", newDiago.GetExtensions()[0].GetPanelHtml(nil))
		assert.Equal(t, "<div>Panel2</div>", newDiago.GetExtensions()[1].GetPanelHtml(nil))
	})

	t.Run("Test BeforeNext on Extension", func(t *testing.T) {
		newDiago := diago.NewDiago()

		mockExtension := &MockDiagoExtension{}

		newDiago.AddExtension(mockExtension)

		mockExtension.BeforeNext(nil)

		assert.True(t, mockExtension.BeforeCalled)
	})

	t.Run("Test AfterNext on Extension", func(t *testing.T) {
		newDiago := diago.NewDiago()

		mockExtension := &MockDiagoExtension{}

		newDiago.AddExtension(mockExtension)

		mockExtension.AfterNext(nil)

		assert.True(t, mockExtension.AfterCalled)
	})

	t.Run("Test GetPanelHtml from Extension", func(t *testing.T) {
		newDiago := diago.NewDiago()

		mockExtension := &MockDiagoExtension{PanelHtml: "<div>Panel</div>"}

		newDiago.AddExtension(mockExtension)

		assert.Equal(t, "<div>Panel</div>", newDiago.GetExtensions()[0].GetPanelHtml(nil))
	})

	t.Run("Test GetHtml from Extension", func(t *testing.T) {
		newDiago := diago.NewDiago()

		mockExtension := &MockDiagoExtension{Html: "<div>Content</div>"}

		newDiago.AddExtension(mockExtension)

		assert.Equal(t, "<div>Content</div>", newDiago.GetExtensions()[0].GetHtml(nil))
	})

	t.Run("Test GetJSHtml from Extension", func(t *testing.T) {
		newDiago := diago.NewDiago()

		mockExtension := &MockDiagoExtension{JSHtml: "<script>console.log('test');</script>"}

		newDiago.AddExtension(mockExtension)

		assert.Equal(t, "<script>console.log('test');</script>", newDiago.GetExtensions()[0].GetJSHtml(nil))
	})

	t.Run("Test ContainsMIME", func(t *testing.T) {
		newDiago := diago.NewDiago()

		assert.True(t, newDiago.ContainsMIME(diago.ContentType_PLAIN))
	})

	t.Run("Test ContainsMIME false", func(t *testing.T) {
		newDiago := diago.NewDiago()

		assert.False(t, newDiago.ContainsMIME("application/json; charset=test"))
	})

	t.Run("Test ContainsMIME false not contain charset", func(t *testing.T) {
		newDiago := diago.NewDiago()

		assert.False(t, newDiago.ContainsMIME("application/json"))
	})
}
