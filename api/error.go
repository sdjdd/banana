package main

import "github.com/labstack/echo/v4"

type httpErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	// server error
	errForbidden = echo.NewHTTPError(403, "forbidden")
	errNotExists = echo.NewHTTPError(404, "no such file of directiry")
	errInternal  = echo.NewHTTPError(500, "internal server error")

	// fs error
	errNoFilename    = echo.NewHTTPError(400, httpErr{101, "filename must be provided"})
	errBaseNotExists = echo.NewHTTPError(400, httpErr{102, "base directory not exists"})
	errIsNotDir      = echo.NewHTTPError(400, httpErr{103, "is not directory"})
	errExists        = echo.NewHTTPError(400, httpErr{104, "already exists"})
	errInsufficient  = echo.NewHTTPError(400, httpErr{105, "space is insufficient"})
)
