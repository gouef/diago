package extensions

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gouef/diago"
	"log"
	"time"
)

type LatencyData struct {
	Latency string
}

type LatencyExtension struct {
	startTime        time.Time
	latency          time.Duration
	PanelGenerator   diago.PanelGenerator
	TemplateProvider diago.TemplateProvider
}

type DefaultLatencyTemplateProvider struct{}

func (p *DefaultLatencyTemplateProvider) GetTemplate() string {
	return diago.GetDiagoLatencyPanelTemplate()
}

func NewDefaultTemplateProvider() *DefaultLatencyTemplateProvider {
	return &DefaultLatencyTemplateProvider{}
}

func NewLatencyExtension() *LatencyExtension {
	generator := diago.NewDefaultPanelGenerator()
	tmpProvider := NewDefaultTemplateProvider()
	return &LatencyExtension{
		PanelGenerator:   generator,
		TemplateProvider: tmpProvider,
	}
}

func (e *LatencyExtension) SetTemplateProvider(provider diago.TemplateProvider) {
	e.TemplateProvider = provider
}

func (e *LatencyExtension) GetTemplateProvider() diago.TemplateProvider {
	return e.TemplateProvider
}

func (e *LatencyExtension) SetPanelGenerator(generator diago.PanelGenerator) {
	e.PanelGenerator = generator
}

func (e *LatencyExtension) GetPanelGenerator() diago.PanelGenerator {
	return e.PanelGenerator
}

func (e *LatencyExtension) GetLatency() time.Duration {
	return e.latency
}

func (e *LatencyExtension) SetLatency(latency time.Duration) {
	e.latency = latency
}

func (e *LatencyExtension) GetHtml(c *gin.Context) string {
	return ""
}
func (e *LatencyExtension) GetJSHtml(c *gin.Context) string {
	return ""
}

func (e *LatencyExtension) GetPanelHtml(c *gin.Context) string {

	var formattedLatency string
	formattedLatency = formatLatency(e.latency)

	log.Printf("Time: %s", formattedLatency)

	result, err := e.PanelGenerator.GenerateHTML("diagoLatencyPanel", e.TemplateProvider, LatencyData{Latency: formattedLatency})

	if err != nil {
		log.Printf("Diago Lattency Extension: %s", err.Error())
	}
	return result
}

func (e *LatencyExtension) BeforeNext(c *gin.Context) {
	e.startTime = time.Now()
}
func (e *LatencyExtension) AfterNext(c *gin.Context) {
	e.latency = time.Since(e.startTime)
}

func formatLatency(latency time.Duration) string {
	switch {
	case latency > time.Second:
		return fmt.Sprintf("%.2f s", float64(latency)/float64(time.Second))
	case latency > time.Millisecond:
		return fmt.Sprintf("%.2f ms", float64(latency)/float64(time.Millisecond))
	case latency > time.Microsecond:
		return fmt.Sprintf("%.2f µs", float64(latency)/float64(time.Microsecond))
	default:
		return fmt.Sprintf("%.2f ns", float64(latency)/float64(time.Nanosecond))
	}
}
