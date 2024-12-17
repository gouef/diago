package diago

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
