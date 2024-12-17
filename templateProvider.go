package diago

import (
	"html/template"
	"strings"
)

type TemplateProvider interface {
	GetTemplate() string
}

type PanelGenerator interface {
	GenerateHTML(name string, templateProvider TemplateProvider, data interface{}) (string, error)
}

type DefaultPanelGenerator struct{}

type DefaultTemplateProvider struct{}

func (p DefaultTemplateProvider) GetTemplate() string {
	return GetDiagoPanelTemplate()
}

func (d *DefaultPanelGenerator) GenerateHTML(name string, templateProvider TemplateProvider, data interface{}) (string, error) {
	templateContent := templateProvider.GetTemplate()
	tpl, err := template.New(name).Parse(templateContent)
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

func NewDefaultPanelGenerator() *DefaultPanelGenerator {
	return &DefaultPanelGenerator{}
}

func NewDefaultTemplateProvider() *DefaultTemplateProvider {
	return &DefaultTemplateProvider{}
}
