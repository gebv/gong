package server

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"net/http"
)

func SessionMiddleware() echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rq := c.Request().(*standard.Request)

			session, err := CookieStore.Get(rq.Request, "session")

			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}

			c.Set("session", session)

			// return c.NoContent(http.StatusNoContent)

			return next(c)
		}
	}
}
