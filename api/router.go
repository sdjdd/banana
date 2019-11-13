package main

import "github.com/labstack/echo/v4"

func loadRouters(e *echo.Echo) {
	e.GET("/info", handleGetInfo)
	e.POST("/info", handleForbidden)
	e.DELETE("/info", handleForbidden)

	e.GET("/*", handleGetFile)
	e.POST("/*", handlePostFile)
	e.DELETE("/*", handleDelFile)
}

func handleForbidden(c echo.Context) error {
	return errForbidden
}
