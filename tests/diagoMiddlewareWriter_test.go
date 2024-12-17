package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockResponseWriter simuluje HTTP ResponseWriter pro testování.
type MockResponseWriter struct {
	mock.Mock
	StatusCode int
	Headers    http.Header
}

func (m *MockResponseWriter) Header() http.Header {
	if m.Headers == nil {
		m.Headers = make(http.Header)
	}
	return m.Headers
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.StatusCode = statusCode
	m.Called(statusCode) // Očekávané volání
}

func (m *MockResponseWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

func TestResponseWriter(t *testing.T) {
	mockWriter := new(MockResponseWriter)

	mockWriter.On("WriteHeader", 200).Once()

	req := httptest.NewRequest("GET", "/", nil)

	handler := gin.Default()
	handler.GET("/", func(c *gin.Context) {
		c.Status(200) // Očekáváme odpověď s kódem 200
	})

	handler.ServeHTTP(mockWriter, req)

	assert.Equal(t, 200, mockWriter.StatusCode)

	mockWriter.AssertExpectations(t)
}
