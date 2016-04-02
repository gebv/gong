package server

import (
    "github.com/labstack/echo"
    // "net/http"
    "html/template"
    // "store"
    "widgets"
)

func SettingsEntryPoint() echo.HandlerFunc {
	return func(c echo.Context) error {
        
        templateString, err := settingsApp.String("index.html")
        
        if err != nil {
            
            return err
        }
        
        tpl, _ := template.New("index").Parse(templateString)
        tpl.Execute(c.Response().Writer(), map[string]string{"Message": "Hello, world!"})
        
        return nil
	}
}

func WebAppEntryPoint() echo.HandlerFunc {
	return func(c echo.Context) error {
        var ClassiferExtId = "pages"
        var ExtId = c.Request().URL().Path()
        
        context := widgets.NewContext()
        
        return context.RenderWidget(c.Response().Writer(), ClassiferExtId, ExtId);
	}
}