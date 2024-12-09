package latency

import (
	_ "embed"
)

//go:embed diago_latency_panel.gohtml
var DiagoLatencyPanelTemplate string

func GetDiagoLatencyPanelTemplate() string {
	return DiagoLatencyPanelTemplate
}
