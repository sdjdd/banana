package main

import "github.com/labstack/echo/v4"

type httpErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	// server error
	errForbidden = echo.NewHTTPError(403, "Forbidden")
	errNotExists = echo.NewHTTPError(404, "No such file of directiry")
	errInternal  = echo.NewHTTPError(500, "Internal server error")

	// fs error
	errNoFilename    = echo.NewHTTPError(400, httpErr{1, "Filename must be provided"})
	errBaseNotExists = echo.NewHTTPError(400, httpErr{2, "Base directory not exists"})
	errIsNotDir      = echo.NewHTTPError(400, httpErr{3, "Base is not directory"})
	errExists        = echo.NewHTTPError(400, httpErr{4, "Already exists"})
	errInsufficient  = echo.NewHTTPError(400, httpErr{5, "Space is insufficient"})
)
