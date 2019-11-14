package main

import "github.com/labstack/echo/v4"

type httpErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	// server error
	errInternal  = echo.NewHTTPError(500, httpErr{100, "internal server error"})
	errForbidden = echo.NewHTTPError(403, httpErr{101, "forbidden"})

	// fs error
	errNoFilename    = echo.NewHTTPError(400, httpErr{200, "filename must be provided"})
	errBaseNotExists = echo.NewHTTPError(400, httpErr{201, "base directory not exists"})
	errNotExists     = echo.NewHTTPError(404, httpErr{202, "no such file of directiry"})
	errIsNotDir      = echo.NewHTTPError(400, httpErr{203, "is not directory"})
	errExists        = echo.NewHTTPError(400, httpErr{204, "already exists"})
	errInsufficient  = echo.NewHTTPError(400, httpErr{205, "space is insufficient"})
)
