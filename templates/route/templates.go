package templates

import (
	_ "embed"
)

//go:embed diago_route_panel.gohtml
var DiagoRoutePanelTemplate string

//go:embed diago_route_panel_js.gohtml
var DiagoRoutePanelJSTemplate string

//go:embed diago_route_panel_popup.gohtml
var DiagoRoutePanelPopupTemplate string

func GetDiagoRoutePanelTemplate() string {
	return DiagoRoutePanelTemplate
}
func GetDiagoRoutePanelJSTemplate() string {
	return DiagoRoutePanelJSTemplate
}
func GetDiagoRoutePanelPopupTemplate() string {
	return DiagoRoutePanelPopupTemplate
}
