package server

import (
	"github.com/GeertJohan/go.rice"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
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

var emptyHandler = func(c echo.Context) error {
	_, err := c.Response().Write([]byte{})
	return err
}

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
	e.Post("/*", WebAppEntryPoint())
	e.Get("/favicon.ico", echo.HandlerFunc(emptyHandler))

	gSettings := e.Group("/@settings")

	gSettings.Get("/*", SettingsEntryPoint())
	gSettings.Get("/static/*", standard.WrapHandler(http.StripPrefix("/@settings/static/", assetHandler)))

	g := gSettings.Group("/api/v1")

	gBuckets := g.Group("/buckets")

	// // Routes classifers
	gBuckets.Get("/:bucket_id", getBucket())
	gBuckets.Get("/search", searchBuckets())
	gBuckets.Put("/:bucket_id", updateBucket())
	gBuckets.Delete("/:bucket_id", deleteBucket())
	gBuckets.Post("", createBucket())

	gItems := g.Group("/files")

	// // Routes items
	gItems.Get("/:file_id", getFile())
	gItems.Get("/search", searchFiles())
	gItems.Put("/:file_id", updateFile())
	gItems.Delete("/:file_id", deleteFile())
	gItems.Post("", createFile())

	//

	// Start server
	e.Run(standard.New(utils.Cfg.Api.Bind))
}
