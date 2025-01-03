package tests

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gouef/diago"
	"github.com/gouef/diago/extensions"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDiagoLatencyExtension(t *testing.T) {
	latencyExtension := extensions.NewLatencyExtension()

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

func TestDiagoLatencyExtension_GetLatency(t *testing.T) {
	latencyExtension := extensions.NewLatencyExtension()
	latencyExtension.SetLatency(2 * time.Second)

	latency := latencyExtension.GetLatency()
	assert.Equal(t, latency, 2*time.Second, "Latency should be equal")
}

func TestDiagoLatencyExtension_GetHtml(t *testing.T) {
	latencyExtension := extensions.NewLatencyExtension()

	html := latencyExtension.GetHtml(nil)
	assert.Empty(t, html, "GetHtml should return an empty string")
}

func TestDiagoLatencyExtension_GetJSHtml(t *testing.T) {
	latencyExtension := extensions.NewLatencyExtension()

	jsHtml := latencyExtension.GetJSHtml(nil)
	assert.Empty(t, jsHtml, "GetJSHtml should return an empty string")
}

func TestDiagoLatencyExtension_GetPanelHtml_Seconds(t *testing.T) {
	latencyExtension := extensions.NewLatencyExtension()
	latencyExtension.SetLatency(2 * time.Second)

	panelHtml := latencyExtension.GetPanelHtml(nil)
	assert.Contains(t, panelHtml, "2.00 s", "Panel HTML should contain latency in seconds")
}

func TestDiagoLatencyExtension_GetPanelHtml_Milliseconds(t *testing.T) {
	latencyExtension := extensions.NewLatencyExtension()
	latencyExtension.SetLatency(150 * time.Millisecond)

	panelHtml := latencyExtension.GetPanelHtml(nil)
	assert.Contains(t, panelHtml, "150.00 ms", "Panel HTML should contain latency in milliseconds")
}

func TestDiagoLatencyExtension_GetPanelHtml_Microseconds(t *testing.T) {
	latencyExtension := extensions.NewLatencyExtension()
	latencyExtension.SetLatency(500 * time.Microsecond)

	panelHtml := latencyExtension.GetPanelHtml(nil)
	assert.Contains(t, panelHtml, "500.00 µs", "Panel HTML should contain latency in microseconds")
}

func TestDiagoLatencyExtension_GetPanelHtml_Nanoseconds(t *testing.T) {
	latencyExtension := extensions.NewLatencyExtension()
	latencyExtension.SetLatency(500 * time.Nanosecond)

	panelHtml := latencyExtension.GetPanelHtml(nil)
	assert.Contains(t, panelHtml, "500.00 ns", "Panel HTML should contain latency in nanoseconds")
}

type mockDiagoPanelGeneratorWithError struct{}

func (m *mockDiagoPanelGeneratorWithError) GenerateHTML(name string, templateProvider diago.TemplateProvider, data interface{}) (string, error) {
	return "", errors.New("mock error generating HTML")
}

func TestDiagoLatencyExtension_GetPanelHtml_ErrorHandling(t *testing.T) {
	latencyExtension := extensions.NewLatencyExtension()

	gen := &mockDiagoPanelGeneratorWithError{}
	latencyExtension.PanelGenerator = gen

	latencyExtension.SetLatency(500 * time.Millisecond)

	var logOutput string
	log.SetOutput(&logWriter{&logOutput})

	panelHtml := latencyExtension.GetPanelHtml(nil)
	assert.Empty(t, panelHtml, "Panel HTML should be empty when there's an error")

	assert.Contains(t, logOutput, "Diago Lattency Extension: mock error generating HTML", "Error message should be logged")
}

type mockTemplateProviderWithParseError struct{}

func (m *mockTemplateProviderWithParseError) GetTemplate() string {
	return "{{ .Latencys }}"
}

type mockTemplateProviderWithExecuteError struct{}

func (m *mockTemplateProviderWithExecuteError) GetTemplate() string {
	return `{{ .NonExistentField }}`
}

func TestGenerateDiagoPanelHTML_TemplateParseError(t *testing.T) {
	mockProvider := &mockTemplateProviderWithParseError{}

	latencyExtension := extensions.NewLatencyExtension()

	result, err := latencyExtension.PanelGenerator.GenerateHTML("error", mockProvider, extensions.LatencyData{Latency: "500 ms"})

	assert.Error(t, err, "Expected error while parsing template")
	assert.Empty(t, result, "Expected empty result when parsing fails")
}

type mockInvalidTemplateProvider struct{}

func (m *mockInvalidTemplateProvider) GetDiagoLatencyPanelTemplate() string {
	return "{{ .InvalidField"
}

func TestGenerateDiagoPanelHTML_TemplateExecuteError(t *testing.T) {
	mockProvider := &mockTemplateProviderWithExecuteError{}

	latencyExtension := extensions.NewLatencyExtension()

	result, err := latencyExtension.PanelGenerator.GenerateHTML("test", mockProvider, extensions.LatencyData{Latency: "500 ms"})

	assert.Error(t, err, "Expected error while executing template")
	assert.Empty(t, result, "Expected empty result when execution fails")
}

func TestDiagoLatencyExtension_TemplateProvider(t *testing.T) {
	latencyExtension := extensions.NewLatencyExtension()

	mockProvider := &mockTemplateProviderWithParseError{}

	latencyExtension.SetTemplateProvider(mockProvider)

	assert.Equal(t, mockProvider, latencyExtension.GetTemplateProvider(), "TemplateProvider should be set correctly")
}

func TestDiagoLatencyExtension_PanelGenerator(t *testing.T) {
	latencyExtension := extensions.NewLatencyExtension()

	mockProvider := &MockPanelGenerator{}

	latencyExtension.SetPanelGenerator(mockProvider)

	assert.Equal(t, mockProvider, latencyExtension.GetPanelGenerator(), "PanelGenerator should be set correctly")
}

type logWriter struct {
	output *string
}

func (lw *logWriter) Write(p []byte) (n int, err error) {
	*lw.output = string(p)
	return len(p), nil
}
