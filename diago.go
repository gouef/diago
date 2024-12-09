package diago

import "github.com/gin-gonic/gin"

type Diago struct {
	Extensions []DiagoExtension
}

func NewDiago() *Diago {
	return &Diago{}
}

func (d *Diago) GetExtensions() []DiagoExtension {
	return d.Extensions
}

func (d *Diago) AddExtension(extension DiagoExtension) *Diago {
	d.Extensions = append(d.Extensions, extension)
	return d
}

type DiagoExtension interface {
	GetPanelHtml(c *gin.Context) string
	GetHtml(c *gin.Context) string
	GetJSHtml(c *gin.Context) string
	BeforeNext(c *gin.Context)
	AfterNext(c *gin.Context)
}
