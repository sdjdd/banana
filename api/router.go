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

	// privileges
	canDownload := userPrivilege{Download: true}
	canUpload := userPrivilege{Upload: true}
	canDelete := userPrivilege{Delete: true}
	canMove := userPrivilege{Upload: true, Delete: true}

	middlewareBasicAuth := middleware.BasicAuth(auth)

	e.GET("/api/whoami", handleWhoAmI)

	fs := e.Group("/fs", middlewareBasicAuth, middlewareParsePath)
	{
		fs.GET("/*", handleGetFile, verifyPrivilege(canDownload))
		fs.POST("/*", handlePostFile, verifyPrivilege(canUpload))
		fs.DELETE("/*", handleDelFile, verifyPrivilege(canDelete))
	}

	api := e.Group("/api", middlewareBasicAuth)
	{
		api.GET("/info", handleGetInfo)
		api.POST("/mv", handleMoveFile, verifyPrivilege(canMove))
	}
}

func auth(username, password string, c echo.Context) (bool, error) {
	user, ok := users.verify(username, password)
	if !ok {
		return false, errForbidden
	}
	c.Set("user", user)
	return true, nil
}

func middlewareParsePath(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		path := strings.TrimLeft(c.Request().URL.Path, "/fs")
		c.Set("path", filepath.Join(diskRoot, path))
		return next(c)
	}
}

func verifyPrivilege(mustCan userPrivilege) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*user)
			if mustCan.Download && !user.Privilege.Download ||
				mustCan.Upload && !user.Privilege.Upload ||
				mustCan.Delete && !user.Privilege.Delete {
				return errForbidden
			}
			return next(c)
		}
	}
}
