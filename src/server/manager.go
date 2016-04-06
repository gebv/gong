package server

import (
    "github.com/labstack/echo"
    "github.com/golang/glog"
    "net/http"
    "store"
    "strings"
    "widgets"
    "bytes"
)

type ResponseError struct {
    Message string
}

func searchItem() echo.HandlerFunc {
    return func(c echo.Context) error {
        var query = c.QueryParam("q");
        var classifer_id = strings.TrimSpace(c.QueryParam("classifer_id"));
        
        if len(classifer_id) == 0 {
            
            return c.JSON(http.StatusBadRequest, ResponseError{"classifer_id is required"})
        }
        
        items, err := store.ItemSearchQuery(query, classifer_id)
        
        // TODO: тестовую генерацию вынести на уровень выше вынести на уровень ниже
        
        for _, widget := range items {
            ctx := widgets.NewContext(c)
            if err := ctx.RenderWidget(bytes.NewBufferString(""), classifer_id, widget.ExtId); err != nil {
                 widget.Props["_BuildError"] = err.Error()
            }
            
            widget.Props["_BuildTraceWidgets"] = ctx.TraceWidgets
        }
        
        if err != nil {
            return err
        }
        
        return c.JSON(http.StatusOK, items)
    }
}

func searchClassifer() echo.HandlerFunc {
    return func(c echo.Context) error {
        var query = c.QueryParam("q");
        
        items, err := store.ItemSearchQuery(query, store.CLASSIFERS)
        
        if err != nil {
            return err
        }
        
        return c.JSON(http.StatusOK, items)
    }
}

func createClassifer() echo.HandlerFunc {
	return func(c echo.Context) error {
        dto := store.NewCreateItemDTO()
         
        if err := c.Bind(dto); err != nil {
            return err
        }
        
        model, err := store.CreateClassifer(dto)
		
        if err != nil {
            return err
        }
        
        return c.JSON(http.StatusCreated, model)
	}
}

func createItem() echo.HandlerFunc {
	return func(c echo.Context) error {
        dto := store.NewCreateItemDTO()
        
        var classifer_id = strings.TrimSpace(c.QueryParam("classifer_id"));
        
        if len(classifer_id) == 0 {
            
            return c.JSON(http.StatusBadRequest, ResponseError{"classifer_id is required"})
        }
         
        if err := c.Bind(dto); err != nil {
            return err
        }
        
        model, err := store.CreateItem(dto, classifer_id)
		
        if err != nil {
            return err
        }
        
        return c.JSON(http.StatusCreated, model)
	}
}


func getItem() echo.HandlerFunc {
	return func(c echo.Context) error {
        modelId := c.Param("model_id")
        
        model, err := store.GetItem(modelId)
        
        if err != nil {
            glog.Infof("%v", err)
           return err 
        }
        
        return c.JSON(http.StatusOK, model)
	}
}

func updateItem() echo.HandlerFunc {
	return func(c echo.Context) error {
        modelId := c.Param("model_id")
         dto := store.NewUpdateItemDTO()
         
        if err := c.Bind(dto); err != nil {
            return err
        }
        
        model, err := store.UpdateItem(modelId, dto, false)
        
        if err != nil {
            
           return err 
        }
        
        return c.JSON(http.StatusOK, model)
	}
}

func deleteClassifer() echo.HandlerFunc {
	return func(c echo.Context) error {
        modelId := c.Param("model_id")
        dto := store.NewUpdateItemDTO()
        
        model, err := store.UpdateItem(modelId, dto, true)
        
        if err != nil {
            
           return err 
        }
        
        return c.JSON(http.StatusOK, model)
	}
}