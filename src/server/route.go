package server

import (
	"github.com/GeertJohan/go.rice"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	// "github.com/labstack/echo/engine/fasthttp"

	// "github.com/labstack/echo/middleware"
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

// Echo instance
var e = echo.New()

func RunServer() {

	// Middleware
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	// e.Use(middleware.Gzip())

	// assetHandler := http.FileServer(rice.MustFindBox("web-static").HTTPBox())
	settingsApp = rice.MustFindBox("settings-app")

	e.Get("/", WebAppEntryPointNew())
	e.Get("/*", WebAppEntryPointNew())
	e.Post("/*", WebAppEntryPointNew())

	e.Static("/favicon.ico", "./src/server/web-static/favicon.ico")
	e.Static("/@settings/static", "./src/server/web-static") // TODO: Вынести в корневой каталог

	gSettings := e.Group("/@settings")

	gSettings.Get("/*", SettingsEntryPoint())
	// gSettings.Get("/static/*", standard.WrapHandler(http.StripPrefix("/@settings/static/", assetHandler)))

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

	e.Router()
	e.Run(standard.New(utils.Cfg.Api.Bind))
}
