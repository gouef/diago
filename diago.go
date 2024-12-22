package diago

import (
	"github.com/gouef/utils"
	"strings"
)

type ContentType struct {
	Types    []string
	Charsets []string
}

const (
	ContentType_HTML  = "text/html"
	ContentType_PLAIN = "text/plain"

	Charset_UTF8 = "utf-8"
	Charset_ALL  = "*"
)

type Diago struct {
	Extensions          []Extension
	TemplateProvider    TemplateProvider
	PanelGenerator      PanelGenerator
	AllowedContentTypes ContentType
}

func NewDiago() *Diago {

	return &Diago{
		TemplateProvider: NewDefaultTemplateProvider(),
		PanelGenerator:   NewDefaultPanelGenerator(),
		AllowedContentTypes: ContentType{
			Types: []string{
				ContentType_HTML,
				ContentType_PLAIN,
			},
			Charsets: []string{
				Charset_ALL,
			},
		},
	}
}

func (d *Diago) SetAllowedContentTypes(contentType ContentType) *Diago {
	d.AllowedContentTypes = contentType
	return d
}

func (d *Diago) AddContentType(typeString string) *Diago {
	d.AllowedContentTypes.Types = append(d.AllowedContentTypes.Types, typeString)
	return d
}

func (d *Diago) AddContentCharset(charset string) *Diago {
	d.AllowedContentTypes.Types = append(d.AllowedContentTypes.Charsets, charset)
	return d
}

func (d *Diago) ContainsMIME(header string) bool {
	parts := strings.Split(header, ";")
	if len(parts) < 1 {
		return false
	}
	contentType := strings.TrimSpace(parts[0])

	charset := ""
	if len(parts) > 1 {
		for _, part := range parts[1:] {
			if strings.HasPrefix(strings.TrimSpace(part), "charset=") {
				part = strings.TrimSpace(part)
				charset = strings.TrimSpace(strings.TrimPrefix(part, "charset="))
				break
			}
		}
	}

	return utils.InListArray([]string{"*", charset}, d.AllowedContentTypes.Charsets) && utils.InArray(contentType, d.AllowedContentTypes.Types)
}

func (d *Diago) GetExtensions() []Extension {
	return d.Extensions
}

func (d *Diago) AddExtension(extension Extension) *Diago {
	d.Extensions = append(d.Extensions, extension)
	return d
}
