package main

import (
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func loadRouters(e *echo.Echo) {
	// e.GET("/", func(c echo.Context) error {
	// 	return c.Redirect(301, "/ui/")
	// })

	middlewareBasicAuth := middleware.BasicAuth(auth)

	fs := e.Group("/fs", middlewareBasicAuth, middlewareParsePath)
	{
		fs.GET("/*", handleGetFile)
		fs.POST("/*", handlePostFile)
		fs.DELETE("/*", handleDelFile)
	}

	api := e.Group("/api", middlewareBasicAuth)
	{
		api.GET("/info", handleGetInfo)
		api.POST("/mv", handleMoveFile)
	}

	e.GET("/api/auth/verify", handleVerifyAuth)
}

func auth(username, password string, c echo.Context) (bool, error) {
	user, ok := users.verify(username, password)
	if !ok {
		return false, nil
	}
	c.Set("user", user)
	return true, nil
}

func handleForbidden(c echo.Context) error {
	return errForbidden
}

func middlewareParsePath(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		path := strings.TrimLeft(c.Request().URL.Path, "/fs")
		path, err = filepath.Abs(filepath.Join(conf.Root, path))
		if err != nil {
			c.Logger().Errorf("get absolute path: %s", err)
			return errInternal
		}
		c.Set("path", path)
		return next(c)
	}
}

func handleVerifyAuth(c echo.Context) error {
	status := 403
	usr, pwd, ok := c.Request().BasicAuth()
	if ok {
		if _, ok := users.verify(usr, pwd); ok {
			status = 200
		}
	}
	return c.NoContent(status)
}
