package main

import (
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

func loadRouters(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(301, "/ui/")
	})

	e.GET("/info", handleGetInfo)

	fs := e.Group("/fs", middlewareParsePath)
	{
		fs.GET("/*", handleGetFile)
		fs.POST("/*", handlePostFile)
		fs.DELETE("/*", handleDelFile)
	}

	e.POST("/mv", handleMoveFile)
}

func handleForbidden(c echo.Context) error {
	return errForbidden
}

func middlewareParsePath(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		path := strings.TrimLeft(c.Request().URL.Path, "/fs")
		path, err = filepath.Abs(filepath.Join(conf.Root, path))
		if err != nil {
			c.Logger().Errorf("get abs path: %s", err)
			return errInternal
		}
		c.Set("path", path)
		return next(c)
	}
}
