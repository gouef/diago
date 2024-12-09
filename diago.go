package diago

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
