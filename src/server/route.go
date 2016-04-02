package server

import (
    "github.com/GeertJohan/go.rice"
    
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    "github.com/labstack/echo/engine/standard"
    "net/http"
    "utils"
)

// Handler
func hello() echo.HandlerFunc {
    return func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!\n")
    }
}

var settingsApp *rice.Box

func RunServer() {
    // Echo instance
    e := echo.New()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    
    assetHandler := http.FileServer(rice.MustFindBox("web-static").HTTPBox())
    settingsApp = rice.MustFindBox("settings-app")
    
    e.Get("/", WebAppEntryPoint())
    e.Get("/*", WebAppEntryPoint())
    
    gSettings := e.Group("/@settings")
    
    gSettings.Get("/*", SettingsEntryPoint())
    gSettings.Get("/static/*", standard.WrapHandler(http.StripPrefix("/@settings/static/", assetHandler)))
    
    g := gSettings.Group("/api/v1")
    
    gClassifers := g.Group("/resources/classifers")
    
    // Routes classifers
    gClassifers.Get("/:model_id", getItem())
    gClassifers.Get("/search", searchClassifer())
    gClassifers.Put("/:model_id", updateItem())
    gClassifers.Delete("/:model_id", deleteClassifer())
    gClassifers.Post("", createClassifer())
    
    gItems := g.Group("/resources/items")
    
    // Routes items
    gItems.Get("/:model_id", getItem())
    gItems.Get("/search", searchItem())
    gItems.Put("/:model_id", updateItem())
    gItems.Delete("/:model_id", deleteClassifer())
    gItems.Post("", createItem())
    
    
    // 

    // Start server
    e.Run(standard.New(utils.Cfg.Api.Bind))
}