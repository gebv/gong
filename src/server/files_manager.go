package server

import (
	"bytes"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	"net/http"
	"store"
	"strconv"
	// "strings"
	"widgets"
)

// func getBucket() echo.HandlerFunc {
//     return func(c echo.Context) error {
//         return nil
//     }
// }

func getFile() echo.HandlerFunc {
	return func(c echo.Context) error {
		// bucket_id := c.Param("bucket_id")
		bucket_id := c.QueryParam("bucket_id")
		id := c.Param("file_id")

		if !store.IsExistBucket(bucket_id) {

			return c.JSON(http.StatusBadRequest, F(CodeInvalidData, "not found bucket or disabled bucket"))
		}

		file, err := store.NewOrLoadFile(bucket_id, id)

		if err != nil {

			if err == store.ErrNotFound {
				return c.JSON(http.StatusBadRequest, F(CodeNotFound))
			}

			glog.Errorf("file get: get bucket=%v, file=%v, err=%v", bucket_id, id, err)

			return c.JSON(http.StatusBadRequest, F(CodeUnknown))
		}

		if file.IsNew() {
			return c.JSON(http.StatusBadRequest, F(CodeNotFound, "file is not found or another reason"))
		}

		return c.JSON(http.StatusOK, OK(file))
	}
}

func updateFile() echo.HandlerFunc {
	return func(c echo.Context) error {
		// bucket_id := c.Param("bucket_id")
		bucket_id := c.QueryParam("bucket_id")
		id := c.Param("file_id")

		if !store.IsExistBucket(bucket_id) {

			return c.JSON(http.StatusBadRequest, F(CodeInvalidData, "not found bucket or disabled bucket"))
		}

		file, err := store.NewOrLoadFile(bucket_id, id)

		if err != nil {
			if err == store.ErrNotFound {
				return c.JSON(http.StatusBadRequest, F(CodeNotFound))
			}

			glog.Errorf("file update: get bucket=%v, file=%v, err=%v", bucket_id, id, err)

			return c.JSON(http.StatusBadRequest, F(CodeUnknown))
		}

		if file.IsNew() {
			return c.JSON(http.StatusBadRequest, F(CodeNotFound, "file is not found or another reason"))
		}

		dto := store.NewUpdateFileDTO()

		if err := c.Bind(dto); err != nil {
			glog.Errorf("file update: parse data bucket=%v, file=%v, err=%v", bucket_id, id, err)

			return c.JSON(http.StatusBadRequest, F(CodeInvalidData))
		}

		if err := file.TransformFrom(dto); err != nil {
			glog.Errorf("file update: transform data bucket=%v, file=%v, err=%v", bucket_id, id, err)

			return err
		}

		if err := store.UpdateFile(file); err != nil {
			glog.Errorf("file update: upsert bucket=%v, file=%v, err=%v", bucket_id, id, err)

			return c.JSON(http.StatusBadRequest, F(CodeUnknown))
		}

		return c.JSON(http.StatusOK, OK(file))
	}
}

func createFile() echo.HandlerFunc {
	return func(c echo.Context) error {
		// bucket_id := c.Param("bucket_id")
		bucket_id := c.QueryParam("bucket_id")

		if !store.IsExistBucket(bucket_id) {

			return c.JSON(http.StatusBadRequest, F(CodeInvalidData, "not found bucket or disabled bucket"))
		}

		dto := store.NewCreateFileDTO()

		if err := c.Bind(dto); err != nil {
			glog.Errorf("file create: parse data err=%v", err)

			return c.JSON(http.StatusBadRequest, F(CodeInvalidData))
		}

		file, err := store.NewOrLoadFile(bucket_id, dto.ExtId)

		glog.V(2).Info("find file id=%v, [%v]", dto.ExtId, file)

		if err != nil && err != store.ErrNotFound {
			glog.Errorf("file create: get bucket=%v, file=%v, err=%v", bucket_id, dto.ExtId, err)

			return c.JSON(http.StatusBadRequest, F(CodeUnknown))
		}

		if !file.IsNew() {

			return c.JSON(http.StatusBadRequest, F(CodeExisting))
		}

		if err := file.TransformFrom(dto); err != nil {
			glog.Errorf("file create: transform data bucket=%v, file=%v, err=%v", bucket_id, dto.ExtId, err)

			return err
		}

		glog.Infof("%+v", file)

		if err := store.CreateFile(file); err != nil {
			glog.Errorf("file create: bucket=%v, file=%v, err=%v", bucket_id, dto.ExtId, err)

			return c.JSON(http.StatusBadRequest, F(CodeUnknown))
		}

		return c.JSON(http.StatusOK, OK(file))
	}
}

func deleteFile() echo.HandlerFunc {
	return func(c echo.Context) error {
		// bucket_id := c.Param("bucket_id")
		bucket_id := c.QueryParam("bucket_id")
		id := c.Param("file_id")

		if !store.IsExistBucket(bucket_id) {

			return c.JSON(http.StatusBadRequest, F(CodeInvalidData, "not found bucket or disabled bucket"))
		}

		file, err := store.NewOrLoadFile(bucket_id, id)

		if err != nil {
			if err == store.ErrNotFound {
				return c.JSON(http.StatusBadRequest, F(CodeNotFound))
			}

			glog.Errorf("bucket file: get bucket=%v, file=%v, err=%v", bucket_id, id, err)

			return c.JSON(http.StatusBadRequest, F(CodeUnknown))
		}

		if file.IsNew() {
			return c.JSON(http.StatusBadRequest, F(CodeNotFound, "file is not found or another reason"))
		}

		if err := store.DeleteFile(file); err != nil {
			glog.Errorf("bucket file: bucket=%v, file=%v, err=%v", bucket_id, id, err)

			return c.JSON(http.StatusBadRequest, F(CodeUnknown))
		}

		return c.JSON(http.StatusOK, OK(file))
	}
}

func searchFiles() echo.HandlerFunc {
	return func(c echo.Context) error {
		// bucket_id := c.Param("bucket_id")
		bucket_id := c.QueryParam("bucket_id")

		if !store.IsExistBucket(bucket_id) {

			return c.JSON(http.StatusBadRequest, F(CodeInvalidData, "not found bucket or disabled bucket"))
		}

		bucket, err := store.NewOrLoadBucket(bucket_id)

		if err != nil {
			glog.Errorf("search files: bucket=%v, err=%v", bucket_id, err)

			return c.JSON(http.StatusBadRequest, F(CodeNotFound))
		}

		filter := store.NewSearchFileter()
		filter.Query = c.QueryParam("q")
		filter.Page, _ = strconv.Atoi(c.QueryParam("page"))
		filter.AddCollections(store.CollNameFile)
		filter.AddCollections(bucket.Id.String())
		filter.SetHasEnabled(true) // только не удаленные
		result := store.SearchPerPage(filter)

		for _, widget := range result.Items {
			ctx := widgets.NewC(c)
			ctx.InitSettings()

			if err := ctx.ExecuteFile(bytes.NewBufferString(""), bucket_id, widget.ExtId, true); err != nil {
				widget.Props["_BuildError"] = err.Error()
			}

			widget.Props["_BuildTraceWidgets"] = ctx.GettraceFiles()
		}

		return c.JSON(http.StatusOK, OK(result))
	}
}
