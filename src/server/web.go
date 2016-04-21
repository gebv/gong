package server

import (
	"github.com/golang/glog"
	"github.com/labstack/echo"
	// "net/http"
	"html/template"
	// "store"
	"bytes"
	"strings"
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

func WebAppEntryPointNew() echo.HandlerFunc {
	return func(c echo.Context) error {
		context := widgets.NewC(c)
		buff := bytes.NewBufferString("")

		if err := context.Run(buff); err != nil {
			return err
		}

		// Response status code

		status := context.GetOneGlobalInt("status")

		if status == 0 {
			status = 200
		}

		// Response content type

		glog.Infof(">> %v", context.GetOneSelfString("render"))

		renderType := strings.ToLower(context.GetOneSelfString("render"))

		if len(renderType) == 0 {
			renderType = strings.ToLower(context.Settings.Routing.DefaultRenderMode)
		}

		switch renderType {
		case "html":
			c.HTML(status, buff.String())
		case "string":
			c.String(status, buff.String())
		case "json":
			c.JSON(status, context.J)
		}

		return nil
	}
}

// func WebAppEntryPoint() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		var key = c.Request().URL().Path()

// 		glog.Infof("WebAppEntryPoint -> %v", len(c.Echo().Routes()))

// 		c.Echo().Router().Add("GET", c.Request().URL().QueryParam("path"), echo.HandlerFunc(func(c echo.Context) error {
// 			glog.Infof(c.Param("id"))

// 			glog.Infof("\t\t -> %v", len(c.Echo().Routes()))
// 			return nil
// 		}), c.Echo())

// 		// c.Echo().Router().Add("GET", "/wqd"+c.Request().URL().QueryParam("path"), echo.HandlerFunc(func(c echo.Context) error {
// 		// 	glog.Infof(c.Param("id"))

// 		// 	glog.Infof("%v", len(c.Echo().Routes()))
// 		// 	return nil
// 		// }), c.Echo())

// 		// c = echo.NewContext(c.Request(), c.Response(), c.Echo())

// 		context := widgets.NewContext(c)

// 		if err := context.Execute(c.Response().Writer(), key); err != nil {
// 			glog.Errorf("WebAppEntryPoint: key=%v, err=%v", key, err)
// 			return err
// 		}

// 		return nil
// 	}
// }
