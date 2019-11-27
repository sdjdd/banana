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
	diskSize int64
	diskUsed int64
	diskRoot string

	users = userList{}
)

func main() {
	var confName, staticDir string
	flag.StringVar(&confName, "c", "banana.yml", "config file")
	flag.StringVar(&staticDir, "s", "", "static directory")
	flag.Parse()

	cfg, err := readConfig(confName)
	if err != nil {
		launchErr("read config file:", err)
	}
	if err = users.set(cfg.Users); err != nil {
		launchErr("set users:", err)
	}
	if diskSize, err = parseSize(cfg.Size); err != nil {
		launchErr("parse size:", err)
	}
	if diskUsed, err = getDirSize(cfg.Root); err != nil {
		launchErr("get root size:", err)
	}
	diskRoot = cfg.Root

	e := newEcho(staticDir)
	e.Logger.Info("banana started on ", cfg.Listen)
	launchErr(e.Start(cfg.Listen))
}

func getDirSize(path string) (size int64, err error) {
	err = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err == nil {
			if !info.IsDir() {
				size += info.Size()
			}
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
