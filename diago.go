package diago

import "github.com/gin-gonic/gin"

type Extension interface {
	GetPanelHtml(c *gin.Context) string
	GetHtml(c *gin.Context) string
	GetJSHtml(c *gin.Context) string
	BeforeNext(c *gin.Context)
	AfterNext(c *gin.Context)
}

type Diago struct {
	Extensions []Extension

	TemplateProvider TemplateProvider
	PanelGenerator   PanelGenerator
}

func NewDiago() *Diago {

	return &Diago{
		TemplateProvider: NewDefaultTemplateProvider(),
		PanelGenerator:   NewDefaultPanelGenerator(),
	}
}

func (d *Diago) GetExtensions() []Extension {
	return d.Extensions
}

func (d *Diago) AddExtension(extension Extension) *Diago {
	d.Extensions = append(d.Extensions, extension)
	return d
}
