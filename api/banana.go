package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const version = "v0.0.1-beta"

var (
	used  int64
	conf  config
	users = make(userList)
)

func main() {
	var confName, staticDir string
	flag.StringVar(&confName, "c", "banana.yml", "config file")
	flag.StringVar(&staticDir, "s", "", "static directory")
	flag.Parse()

	var err error
	if err = conf.read(confName); err != nil {
		launchErr("read config file:", err)
	}

	used, err = getDirSize(conf.Root)
	if err != nil {
		launchErr("get root size:", err)
	}

	e := newEcho(staticDir)
	e.Logger.Info("banana started on ", conf.Listen)
	launchErr(e.Start(conf.Listen))
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

func newEcho(staticDir string) (e *echo.Echo) {
	e = echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetHeader("${time_rfc3339} ${level}")
	e.Logger.SetLevel(log.DEBUG)

	loadRouters(e)

	return
}

func launchErr(a ...interface{}) {
	fmt.Println(a...)
	os.Exit(1)
}
