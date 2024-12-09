package extensions

import (
	"github.com/gin-gonic/gin"
	"github.com/gouef/router"
	"html/template"
	"log"
	"strings"
)

type DiagoRouteExtension struct {
	router       *router.Router
	currentRoute string
	data         DiagoRouteData
}

type DiagoRoute struct {
	Actual  bool
	Name    string
	Pattern string
	Method  string
}

type DiagoRouteData struct {
	CurrentRoute string
	Routes       []DiagoRoute
}

func NewDiagoRouteExtension(r *router.Router) *DiagoRouteExtension {
	return &DiagoRouteExtension{router: r}
}

func (e *DiagoRouteExtension) GetHtml(c *gin.Context) string {
	result, err := e.generateDiagoPanelPopupHTML(e.data)

	if err != nil {
		log.Printf("Diago Route Extension: %s", err.Error())
	}
	return result
}
func (e *DiagoRouteExtension) GetJSHtml(c *gin.Context) string {
	result, err := e.generateDiagoPanelJSHTML()

	if err != nil {
		log.Printf("Diago Route Extension: %s", err.Error())
	}
	return result
}
func (e *DiagoRouteExtension) GetPanelHtml(c *gin.Context) string {
	result, err := e.generateDiagoPanelHTML(e.data)

	if err != nil {
		log.Printf("Diago Route Extension: %s", err.Error())
	}
	return result
}

func (e *DiagoRouteExtension) BeforeNext(c *gin.Context) {
}
func (e *DiagoRouteExtension) AfterNext(c *gin.Context) {
	e.currentRoute = c.FullPath()
	var routes []DiagoRoute
	for _, route := range e.router.GetRoutes() {
		routes = append(routes, DiagoRoute{
			Actual:  e.currentRoute == route.GetPattern(),
			Name:    route.GetName(),
			Pattern: route.GetPattern(),
			Method:  route.GetMethod().String(),
		})
	}

	e.data = DiagoRouteData{
		CurrentRoute: e.currentRoute,
		Routes:       routes,
	}
}

func (e *DiagoRouteExtension) generateDiagoPanelHTML(data DiagoRouteData) (string, error) {
	tpl, err := template.New("diagoRoutePanel").ParseFiles("templates/route/diago_route_panel.gohtml")
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

func (e *DiagoRouteExtension) generateDiagoPanelPopupHTML(data DiagoRouteData) (string, error) {
	tpl, err := template.New("diagoRoutePanelPopup").ParseFiles("templates/route/diago_route_panel_popup.gohtml")
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

func (e *DiagoRouteExtension) generateDiagoPanelJSHTML() (string, error) {
	tpl, err := template.New("diagoRoutePanelJS").ParseFiles("templates/route/diago_route_panel_js.gohtml")
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	err = tpl.Execute(&builder, nil)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}
