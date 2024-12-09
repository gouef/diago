package diago

import (
	"bytes"
	_ "embed"
	"github.com/gin-gonic/gin"
	"github.com/gouef/router"
	"html/template"
	"log"
	"strings"
)

type DiagoData struct {
	ExtensionsPanelHtml []template.HTML
	ExtensionsJSHtml    []template.HTML
	ExtensionsHtml      []template.HTML
}

func DiagoMiddleware(r *router.Router, d *Diago) gin.HandlerFunc {
	return func(c *gin.Context) {

		for _, e := range d.GetExtensions() {
			e.BeforeNext(c)
		}

		responseBuffer := &bytes.Buffer{}
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			buffer:         responseBuffer,
		}

		c.Next()

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

			diagoPanelHTML, err := generateDiagoPanelHTML(diagoData)

			if err != nil {
				log.Println("Error generating Diago panel HTML:", err)
				return
			}

			writer.buffer.WriteString(diagoPanelHTML)
		}

		_, err := c.Writer.Write(responseBuffer.Bytes())
		if err != nil {
			log.Println("Error writing response:", err)
		}

		// Logování status kódu
		status := c.Writer.Status()
		log.Printf("Status: %d", status)
	}
}

func generateDiagoPanelHTML(data DiagoData) (string, error) {

	tpl, err := template.New("diagoPanel").ParseFiles(GetDiagoPanelTemplate())
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

type responseWriter struct {
	gin.ResponseWriter
	buffer *bytes.Buffer
}

func (w *responseWriter) Write(data []byte) (int, error) {
	return w.buffer.Write(data)
}

func (w *responseWriter) WriteHeaderNow() {
	if !w.Written() {
		w.ResponseWriter.WriteHeaderNow()
	}
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}
