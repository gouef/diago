package diago

import (
	"bytes"
	_ "embed"
	"github.com/gin-gonic/gin"
	"github.com/gouef/router"
	"html/template"
	"log"
	"net/http"
)

type DiagoData struct {
	ExtensionsPanelHtml []template.HTML
	ExtensionsJSHtml    []template.HTML
	ExtensionsHtml      []template.HTML
}

func DiagoMiddleware(r *router.Router, d *Diago) gin.HandlerFunc {
	return func(c *gin.Context) {

		if r != nil && r.IsRelease() {
			c.Next()
			return
		}

		originalWriter := c.Writer
		for _, e := range d.GetExtensions() {
			e.BeforeNext(c)
		}

		responseBuffer := &bytes.Buffer{}
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			buffer:         responseBuffer,
			statusCode:     http.StatusOK,
		}

		c.Writer = writer

		c.Next()
		c.Copy()
		for _, e := range d.GetExtensions() {
			e.AfterNext(c)
		}

		contentType := writer.Header().Get("Content-Type")

		if contentType == "text/html; charset=utf-8" {
			var extensionsHtml []template.HTML
			var extensionsPanelHtml []template.HTML
			var extensionsJSHtml []template.HTML

			for _, e := range d.GetExtensions() {
				extensionsHtml = append(extensionsHtml, template.HTML(e.GetHtml(c)))
				extensionsPanelHtml = append(extensionsPanelHtml, template.HTML(e.GetPanelHtml(c)))
				extensionsJSHtml = append(extensionsJSHtml, template.HTML(e.GetJSHtml(c)))
			}

			diagoData := DiagoData{
				ExtensionsHtml:      extensionsHtml,
				ExtensionsPanelHtml: extensionsPanelHtml,
				ExtensionsJSHtml:    extensionsJSHtml,
			}

			diagoPanelHTML, err := d.PanelGenerator.GenerateHTML("diagoPanel", d.TemplateProvider, diagoData)

			if err != nil {
				err = c.Error(err)
				c.Status(500)
				c.Writer.WriteHeaderNow()
				diagoPanelHTML = "Error generating Diago panel HTML"
			}

			writer.buffer.WriteString(diagoPanelHTML)
		}

		c.Writer = originalWriter
		_, err := c.Writer.Write(responseBuffer.Bytes())
		if err != nil {
			err = c.Error(err)
			c.Status(500)
			writer.buffer.WriteString("Error generating Diago panel HTML")
		}

		status := c.Writer.Status()
		log.Printf("Status: %d", status)
	}
}

type responseWriter struct {
	gin.ResponseWriter
	buffer     *bytes.Buffer
	statusCode int
}

func (w *responseWriter) Write(data []byte) (int, error) {
	return w.buffer.Write(data)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
