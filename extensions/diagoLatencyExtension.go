package extensions

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gouef/diago"
	"html/template"
	"log"
	"strings"
	"time"
)

type DiagoLatencyExtension struct {
	startTime time.Time
	latency   time.Duration
}

func NewDiagoLatencyExtension() *DiagoLatencyExtension {
	return &DiagoLatencyExtension{}
}

func (e *DiagoLatencyExtension) GetHtml(c *gin.Context) string {
	return ""
}
func (e *DiagoLatencyExtension) GetJSHtml(c *gin.Context) string {
	return ""
}

func (e *DiagoLatencyExtension) GetPanelHtml(c *gin.Context) string {

	var formattedLatency string
	switch {
	case e.latency > time.Second:
		formattedLatency = fmt.Sprintf("%.2f s", float64(e.latency)/float64(time.Second))
	case e.latency > time.Millisecond:
		formattedLatency = fmt.Sprintf("%.2f ms", float64(e.latency)/float64(time.Millisecond))
	case e.latency > time.Microsecond:
		formattedLatency = fmt.Sprintf("%.2f µs", float64(e.latency)/float64(time.Microsecond))
	default:
		formattedLatency = fmt.Sprintf("%.2f ns", float64(e.latency)/float64(time.Nanosecond))
	}

	log.Printf("Time: %s", formattedLatency)

	result, err := e.generateDiagoPanelHTML(struct{ Latency string }{Latency: formattedLatency})

	if err != nil {
		log.Printf("Diago Lattency Extension: %s", err.Error())
	}
	return result
}

func (e *DiagoLatencyExtension) BeforeNext(c *gin.Context) {
	e.startTime = time.Now()
}
func (e *DiagoLatencyExtension) AfterNext(c *gin.Context) {
	e.latency = time.Since(e.startTime)
}

func (e *DiagoLatencyExtension) generateDiagoPanelHTML(data struct {
	Latency string
}) (string, error) {
	tpl, err := template.New("diagoLatencyPanel").Parse(diago.GetDiagoLatencyPanelTemplate())
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	err = tpl.Execute(&builder, data)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}
