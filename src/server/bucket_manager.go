package server

import (
	// "bytes"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	"net/http"
	"store"
	"strconv"
	// "strings"
	// "widgets"
)

// func getBucket() echo.HandlerFunc {
//     return func(c echo.Context) error {
//         return nil
//     }
// }

func getBucket() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("bucket_id")

		file, err := store.NewOrLoadBucket(id)

		if err != nil {

			if err == store.ErrNotFound {
				return c.JSON(http.StatusBadRequest, F(CodeNotFound))
			}

			glog.Errorf("bucket get: get bucket=%v, err=%v", id, err)

			return c.JSON(http.StatusBadRequest, F(CodeUnknown))
		}

		if file.IsNew() {
			return c.JSON(http.StatusBadRequest, F(CodeNotFound, "file is not found or another reason"))
		}

		return c.JSON(http.StatusOK, OK(file))
	}
}

func updateBucket() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("bucket_id")
		file, err := store.NewOrLoadBucket(id)

		if err != nil {
			if err == store.ErrNotFound {
				return c.JSON(http.StatusBadRequest, F(CodeNotFound))
			}

			glog.Errorf("bucket update: get bucket=%v, err=%v", id, err)

			return c.JSON(http.StatusBadRequest, F(CodeUnknown))
		}

		if file.IsNew() {
			return c.JSON(http.StatusBadRequest, F(CodeNotFound, "file is not found or another reason"))
		}

		dto := store.NewUpdateFileDTO()

		if err := c.Bind(dto); err != nil {
			glog.Errorf("bucket update: parse data bucket=%v, err=%v", id, err)

			return c.JSON(http.StatusBadRequest, F(CodeInvalidData))
		}

		if err := file.TransformFrom(dto); err != nil {
			glog.Errorf("bucket update: transform data bucket=%v, err=%v", id, err)

			return err
		}

		if err := store.UpsertFile(file); err != nil {
			glog.Errorf("bucket update: upsert bucket=%v, err=%v", id, err)

			return c.JSON(http.StatusBadRequest, F(CodeUnknown))
		}

		return c.JSON(http.StatusOK, OK(file))
	}
}

func createBucket() echo.HandlerFunc {
	return func(c echo.Context) error {
		dto := store.NewCreateFileDTO()

		if err := c.Bind(dto); err != nil {
			glog.Errorf("bucket create: parse data err=%v", err)

			return c.JSON(http.StatusBadRequest, F(CodeInvalidData))
		}

		file, err := store.NewOrLoadBucket(dto.ExtId)

		if err != nil && err != store.ErrNotFound {
			glog.Errorf("bucket create: get bucket=%v, err=%v", dto.ExtId, err)

			return c.JSON(http.StatusBadRequest, F(CodeUnknown))
		}

		if !file.IsNew() {

			return c.JSON(http.StatusBadRequest, F(CodeExisting))
		}

		if err := file.TransformFrom(dto); err != nil {
			glog.Errorf("bucket create: transform data bucket=%v, err=%v", dto.ExtId, err)

			return err
		}

		glog.Infof("%+v", file)

		if err := store.CreateFile(store.CollNameBucket, file); err != nil {
			glog.Errorf("bucket create: bucket=%v, err=%v", dto.ExtId, err)

			return c.JSON(http.StatusBadRequest, F(CodeUnknown))
		}

		return c.JSON(http.StatusOK, OK(file))
	}
}

func deleteBucket() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("bucket_id")
		file, err := store.NewOrLoadBucket(id)

		if err != nil {
			if err == store.ErrNotFound {
				return c.JSON(http.StatusBadRequest, F(CodeNotFound))
			}

			glog.Errorf("bucket delete: get bucket=%v, err=%v", id, err)

			return c.JSON(http.StatusBadRequest, F(CodeUnknown))
		}

		if file.IsNew() {
			return c.JSON(http.StatusBadRequest, F(CodeNotFound, "file is not found or another reason"))
		}

		if err := store.Delete(file); err != nil {
			glog.Errorf("bucket delete: bucket=%v, err=%v", id, err)

			return c.JSON(http.StatusBadRequest, F(CodeUnknown))
		}

		return c.JSON(http.StatusOK, OK(file))
	}
}

func searchBuckets() echo.HandlerFunc {
	return func(c echo.Context) error {
		filter := store.NewSearchFileter()
		filter.Query = c.QueryParam("query")
		filter.Page, _ = strconv.Atoi(c.QueryParam("page"))
		filter.AddCollections(store.CollNameBucket)
		result := store.SearchPerPage(filter)

		return c.JSON(http.StatusOK, OK(result))
	}
}
