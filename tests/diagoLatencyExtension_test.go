package tests

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gouef/diago/extensions"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockDiagoPanelGenerator struct{}

func (m *mockDiagoPanelGenerator) GenerateDiagoPanelHTML(data struct{ Latency string }) (string, error) {
	return "<div>Mocked HTML: " + data.Latency + "</div>", nil
}

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

func TestDiagoLatencyExtension_GetHtml(t *testing.T) {
	latencyExtension := extensions.NewDiagoLatencyExtension()

	html := latencyExtension.GetHtml(nil)
	assert.Empty(t, html, "GetHtml should return an empty string")
}

func TestDiagoLatencyExtension_GetJSHtml(t *testing.T) {
	latencyExtension := extensions.NewDiagoLatencyExtension()

	jsHtml := latencyExtension.GetJSHtml(nil)
	assert.Empty(t, jsHtml, "GetJSHtml should return an empty string")
}

func TestDiagoLatencyExtension_GetPanelHtml_Seconds(t *testing.T) {
	latencyExtension := extensions.NewDiagoLatencyExtension()
	latencyExtension.SetLatency(2 * time.Second)

	panelHtml := latencyExtension.GetPanelHtml(nil)
	assert.Contains(t, panelHtml, "2.00 s", "Panel HTML should contain latency in seconds")
}

func TestDiagoLatencyExtension_GetPanelHtml_Milliseconds(t *testing.T) {
	latencyExtension := extensions.NewDiagoLatencyExtension()
	latencyExtension.SetLatency(150 * time.Millisecond)

	panelHtml := latencyExtension.GetPanelHtml(nil)
	assert.Contains(t, panelHtml, "150.00 ms", "Panel HTML should contain latency in milliseconds")
}

func TestDiagoLatencyExtension_GetPanelHtml_Microseconds(t *testing.T) {
	latencyExtension := extensions.NewDiagoLatencyExtension()
	latencyExtension.SetLatency(500 * time.Microsecond)

	panelHtml := latencyExtension.GetPanelHtml(nil)
	assert.Contains(t, panelHtml, "500.00 µs", "Panel HTML should contain latency in microseconds")
}

func TestDiagoLatencyExtension_GetPanelHtml_Nanoseconds(t *testing.T) {
	latencyExtension := extensions.NewDiagoLatencyExtension()
	latencyExtension.SetLatency(500 * time.Nanosecond)

	panelHtml := latencyExtension.GetPanelHtml(nil)
	assert.Contains(t, panelHtml, "500.00 ns", "Panel HTML should contain latency in nanoseconds")
}

type mockDiagoPanelGeneratorWithError struct{}

func (m *mockDiagoPanelGeneratorWithError) GenerateDiagoPanelHTML(data struct{ Latency string }) (string, error) {
	return "", errors.New("mock error generating HTML") // Simulace chyby
}

func TestDiagoLatencyExtension_GetPanelHtml_ErrorHandling(t *testing.T) {
	latencyExtension := extensions.NewDiagoLatencyExtension()

	latencyExtension.PanelGenerator = &mockDiagoPanelGeneratorWithError{}

	latencyExtension.SetLatency(500 * time.Millisecond)

	var logOutput string
	log.SetOutput(&logWriter{&logOutput})

	panelHtml := latencyExtension.GetPanelHtml(nil)
	assert.Empty(t, panelHtml, "Panel HTML should be empty when there's an error")

	assert.Contains(t, logOutput, "Diago Lattency Extension: mock error generating HTML", "Error message should be logged")
}

type logWriter struct {
	output *string
}

func (lw *logWriter) Write(p []byte) (n int, err error) {
	*lw.output = string(p)
	return len(p), nil
}
