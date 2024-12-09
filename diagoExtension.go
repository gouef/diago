package diago

import "github.com/gin-gonic/gin"

type DiagoExtension interface {
	GetPanelHtml(c *gin.Context) string
	GetHtml(c *gin.Context) string
	GetJSHtml(c *gin.Context) string
	BeforeNext(c *gin.Context)
	AfterNext(c *gin.Context)
}
