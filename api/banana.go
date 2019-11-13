package main

import (
	"crypto/subtle"
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type fileInfo struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
}

type info struct {
	Version  string   `json:"version"`
	Capacity int64    `json:"capacity"`
	Used     int64    `json:"used"`
	User     userInfo `json:"user"`
}

type userInfo struct {
	Name   string `json:"name"`
	Upload bool   `json:"upload"`
	Delete bool   `json:"delete"`
}

type msg struct {
	Message string `json:"message"`
}

type httpErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const version = "v0.0.1-beta"

var (
	errInternal     = echo.NewHTTPError(500, httpErr{100, "internal server error"})
	errForbidden    = echo.NewHTTPError(403, httpErr{201, "forbidden"})
	errCreateRoot   = echo.NewHTTPError(400, httpErr{202, "cannot create /"})
	errDeleteRoot   = echo.NewHTTPError(400, httpErr{203, "cannot delete /"})
	errNotExists    = echo.NewHTTPError(404, httpErr{204, "no such file of directiry"})
	errIsNotDir     = echo.NewHTTPError(400, httpErr{205, "is not directory"})
	errInsufficient = echo.NewHTTPError(400, httpErr{206, "space is insufficient"})
)

var used int64

func main() {
	e := echo.New()
	e.HideBanner = true
	e.Logger.SetHeader("${time_rfc3339} ${level}")

	loadRouters(e)

	e.Use(middleware.BasicAuthWithConfig(
		middleware.BasicAuthConfig{
			Skipper:   skipBasicAuth,
			Validator: authenticate,
		},
	))

	var err error
	used, err = getDirSize(conf.Root)
	if err != nil {
		log.Fatalln("get root dir size:", err)
	}

	// e.HTTPErrorHandler = func(err error, c echo.Context) {
	// 	c.Logger().Error(err)
	// 	e.DefaultHTTPErrorHandler(err, c)
	// }
	e.Start(conf.Listen)
}

func getDirSize(path string) (size int64, err error) {
	err = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return
}

func skipBasicAuth(c echo.Context) bool {
	return conf.Anonymous.Enable
}

func authenticate(username, password string, c echo.Context) (bool, error) {
	user, ok := conf.Users[username]
	if !ok {
		return false, nil
	}
	if subtle.ConstantTimeCompare([]byte(user.Password), []byte(password)) != 1 {
		return false, nil
	}
	c.Set("username", username)
	c.Set("user", user)
	return true, nil
}
