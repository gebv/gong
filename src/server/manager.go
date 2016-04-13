package server

import (
// "bytes"
// "github.com/golang/glog"
// "github.com/labstack/echo"
// "net/http"
// "store"
// "strconv"
// "strings"
// "widgets"
)

// type ResponseError struct {
// 	Message string
// }

// func searchItem() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		var query = c.QueryParam("q")
// 		var bucket_id = strings.TrimSpace(c.Param("bucket_id"))

// 		if len(bucket_id) == 0 {

// 			return c.JSON(http.StatusBadRequest, ResponseError{"bucket_id is required"})
// 		}

// 		page, _ := strconv.Atoi(c.QueryParam("page"))

// 		filter := store.NewSearchFileter()
// 		filter.SetPage(page)
// 		filter.SetQuery(query)
// 		filter.AddCollections(bucket_id)

// 		items := store.SearchPerPage(filter)

// 		// TODO: тестовую генерацию вынести на уровень выше вынести на уровень ниже

// 		for _, widget := range items.Items {
// 			ctx := widgets.NewContext(c)

// 			// if err := ctx.RenderWidget(bytes.NewBufferString(""), widget.Collections[0], widget.ExtId); err != nil {
// 			// 	widget.Props["_BuildError"] = err.Error()
// 			// }

// 			widget.Props["_BuildTraceWidgets"] = ctx.TraceWidgets
// 		}

// 		return c.JSON(http.StatusOK, items)
// 	}
// }

// func searchClassifer() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		var query = c.QueryParam("q")

// 		page, _ := strconv.Atoi(c.QueryParam("page"))

// 		filter := store.NewSearchFileter()
// 		filter.SetPage(page)
// 		filter.SetQuery(query)
// 		filter.AddCollections(store.CollNameBucket)

// 		items := store.SearchPerPage(filter)

// 		return c.JSON(http.StatusOK, items)
// 	}
// }

// func createClassifer() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		file := store.NewFile()

// 		if err := c.Bind(file); err != nil {
// 			return err
// 		}

// 		file.AddCollections(store.CollNameBucket)

// 		err := store.CreateFile(store.CollNameBucket, file)

// 		if err != nil {
// 			return err
// 		}

// 		return c.JSON(http.StatusCreated, file)
// 	}
// }

// func createItem() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		file := store.NewFile()

// 		var bucket_id = strings.TrimSpace(c.Param("bucket_id"))

// 		if len(bucket_id) == 0 {

// 			return c.JSON(http.StatusBadRequest, ResponseError{"bucket_id is required"})
// 		}

// 		if err := c.Bind(file); err != nil {
// 			return err
// 		}

// 		file.AddCollections(store.CollNameFile)

// 		err := store.CreateFile(bucket_id, file)

// 		if err != nil {
// 			return err
// 		}

// 		return c.JSON(http.StatusCreated, file)
// 	}
// }

// func getItem() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		modelId := c.Param("bucket_id")

// 		file, err := store.FindFileById(modelId)

// 		if err != nil {
// 			glog.Infof("%v", err)
// 			return err
// 		}

// 		return c.JSON(http.StatusOK, file)
// 	}
// }

// func updateItem() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		modelId := c.Param("bucket_id")

// 		file, err := store.FindFileById(modelId)

// 		if err != nil {

// 			return err
// 		}

// 		if err := c.Bind(file); err != nil {
// 			return err
// 		}

// 		if err := store.UpsertFile(file); err != nil {

// 			return err
// 		}

// 		return c.JSON(http.StatusOK, file)
// 	}
// }

// func deleteClassifer() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		modelId := c.Param("model_id")

// 		model, err := store.FindFileById(modelId)

// 		if err != nil {

// 			return err
// 		}

// 		return c.JSON(http.StatusOK, model)
// 	}
// }
