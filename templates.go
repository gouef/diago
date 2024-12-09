package diago

import (
	_ "embed"
)

//go:embed templates/diago_panel.gohtml
var DiagoPanelTemplate string

//go:embed templates/latency/diago_latency_panel.gohtml
var DiagoLatencyPanelTemplate string

func GetDiagoPanelTemplate() string {
	return DiagoPanelTemplate
}

func GetDiagoLatencyPanelTemplate() string {
	return DiagoLatencyPanelTemplate
}
