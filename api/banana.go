package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type msg struct {
	Message string `json:"message"`
}

const version = "v0.0.1-beta"

var binDir string

var used int64

func init() {
	binPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	binDir = filepath.Dir(binPath)
}

func main() {
	if err := conf.read("banana.yml"); err != nil {
		fmt.Println("error: read config file:", err)
		os.Exit(1)
	}

	var err error
	used, err = getDirSize(conf.Root)
	if err != nil {
		log.Fatalln("get root dir size:", err)
	}

	// e.HTTPErrorHandler = func(err error, c echo.Context) {
	// 	c.Logger().Error(err)
	// 	e.DefaultHTTPErrorHandler(err, c)
	// }
	e := newEcho()
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

func newEcho() (e *echo.Echo) {
	e = echo.New()
	e.HideBanner = true
	e.Logger.SetHeader("${time_rfc3339} ${level}")

	search := []string{"ui", filepath.Join(binDir, "ui")}
	for _, path := range search {
		if _, err := os.Stat(path); err == nil {
			e.Static("/ui", path)
			break
		}
	}

	loadRouters(e)

	return
}
