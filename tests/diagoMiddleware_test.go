package tests

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gouef/diago"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"html/template"
	"log"
	"net"
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

type MockExtension struct {
	mock.Mock
}

func (m *MockExtension) BeforeNext(c *gin.Context) {
	m.Called(c)
}

func (m *MockExtension) AfterNext(c *gin.Context) {
	m.Called(c)
}

func (m *MockExtension) GetHtml(c *gin.Context) string {
	args := m.Called(c)
	return args.String(0)
}

func (m *MockExtension) GetPanelHtml(c *gin.Context) string {
	args := m.Called(c)
	return args.String(0)
}

func (m *MockExtension) GetJSHtml(c *gin.Context) string {
	args := m.Called(c)
	return args.String(0)
}

type MockResponseWriter struct {
	mock.Mock
}

// Dynamické volání metod
func (m *MockResponseWriter) Write(data []byte) (int, error) {
	args := m.Called(data)
	return args.Int(0), args.Error(1)
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.Called(statusCode)
}

func (m *MockResponseWriter) Header() http.Header {
	args := m.Called()
	return args.Get(0).(http.Header)
}

// Pro metody, které jsou povinné
func (m *MockResponseWriter) Flush() {
	m.Called()
}

func (m *MockResponseWriter) CloseNotify() <-chan bool {
	args := m.Called()
	return args.Get(0).(<-chan bool)
}

func (m *MockResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	args := m.Called()
	return args.Get(0).(net.Conn), args.Get(1).(*bufio.ReadWriter), args.Error(2)
}

func (m *MockResponseWriter) Size() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockResponseWriter) Status() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockResponseWriter) Written() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockResponseWriter) Pusher() http.Pusher {
	return nil
}

func (m *MockResponseWriter) WriteHeaderNow() {
	m.Called()
}

func (m *MockResponseWriter) WriteString(s string) (int, error) {
	args := m.Called(s)
	return args.Int(0), args.Error(1)
}

func TestDiagoMiddleware_WriteResponse_Error(t *testing.T) {
	mockExtension := new(MockExtension)
	mockExtension.On("BeforeNext", mock.Anything).Return()
	mockExtension.On("AfterNext", mock.Anything).Return()
	mockExtension.On("GetHtml", mock.Anything).Return("<div>Mock Extension HTML</div>")
	mockExtension.On("GetPanelHtml", mock.Anything).Return("<div>Mock Panel HTML</div>")
	mockExtension.On("GetJSHtml", mock.Anything).Return("<script>console.log('Mock JS');</script>")

	r := gin.Default()

	diagoInstance := &diago.Diago{}
	diagoInstance.Extensions = []diago.DiagoExtension{mockExtension}

	mockWriter := new(MockResponseWriter)
	mockWriter.On("Write", mock.Anything).Return(0, errors.New("simulovaná chyba při zápisu"))
	mockWriter.On("Header").Return(http.Header{})
	mockWriter.On("Size").Return(0)
	mockWriter.On("Status").Return(0)
	mockWriter.On("Written").Return(false)

	r.Use(func(c *gin.Context) {
		c.Writer = mockWriter
		diago.DiagoMiddleware(nil, diagoInstance)(c)
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	//mockWriter.AssertExpectations(t)

	var buf bytes.Buffer
	log.SetOutput(&buf)

	r.ServeHTTP(w, req)

	assert.Contains(t, buf.String(), "Error writing response:")

	mockWriter.AssertNumberOfCalls(t, "Write", 4)
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
