package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	msg struct {
		Content string `json:"message"`
	}

	fileInfo struct {
		Name    string    `json:"name"`
		IsDir   bool      `json:"isDir"`
		Size    int64     `json:"size"`
		ModTime time.Time `json:"modTime"`
	}

	serverInfo struct {
		Version string `json:"version"`
		Size    int64  `json:"size"`
		Used    int64  `json:"used"`
	}
)

func handleVerifyAuth(c echo.Context) error {
	if usr, pwd, ok := c.Request().BasicAuth(); ok {
		if _, ok = users.verify(usr, pwd); ok {
			return c.NoContent(200)
		}
	}
	return c.NoContent(403)
}

func handleGetInfo(c echo.Context) error {
	return c.JSON(200, serverInfo{
		Version: version,
		Size:    diskSize,
		Used:    atomic.LoadInt64(&diskUsed),
	})
}

func handleGetFile(c echo.Context) error {
	path := c.Get("path").(string)
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errNotExists
		}
		c.Logger().Error("get file info: ", err)
		return errInternal
	}

	if !info.IsDir() {
		return c.Attachment(path, filepath.Base(path))
	}

	childrenInfo, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errNotExists
		}
		c.Logger().Error("get file info: ", err)
		return errInternal
	}

	result := make([]fileInfo, len(childrenInfo))
	for i, info := range childrenInfo {
		result[i] = fileInfo{
			Name:    info.Name(),
			IsDir:   info.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		}
	}
	return c.JSON(200, result)
}

func handlePostFile(c echo.Context) error {
	path := c.Get("path").(string)
	if path == diskRoot {
		return errNoFilename
	}

	if dirInfo, err := os.Stat(filepath.Dir(path)); err != nil {
		if os.IsNotExist(err) {
			return errBaseNotExists
		}
		c.Logger().Error("get base info: ", err)
		return err
	} else if !dirInfo.IsDir() {
		return errIsNotDir
	}

	if c.QueryParam("dir") == "true" {
		if err := os.Mkdir(path, 0755); err != nil {
			if os.IsExist(err) {
				return errExists
			}
			c.Logger().Error("create directory: ", err)
			return errInternal
		}
		return c.JSON(200, msg{"Directory created successfully"})
	}

	var freeSize int64
	if diskSize > 0 {
		if freeSize = diskSize - atomic.LoadInt64(&diskUsed); freeSize <= 0 {
			return errInsufficient
		}
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return errExists
		}
		c.Logger().Error("create file: ", err)
		return errInternal
	}

	body := c.Request().Body
	defer f.Close()

	var copied int64
	if diskSize > 0 {
		const bufSize int64 = 1 << 21 // 2MB
		for {
			written, err := io.CopyN(f, body, bufSize)
			if copied += written; copied > freeSize {
				f.Close()
				os.Remove(path)
				return errInsufficient
			}
			if err != nil {
				if err == io.EOF {
					atomic.AddInt64(&diskUsed, copied)
					break
				}
				c.Logger().Error("copy file: ", err)
				return errInternal
			}
		}
	} else {
		if copied, err = io.Copy(f, body); err != nil {
			c.Logger().Error("copy file: ", err)
			return errInternal
		}
	}

	return c.JSON(200, msg{"File uploaded successfully"})
}

func handleDelFile(c echo.Context) error {
	path := c.Get("path").(string)
	if path == diskRoot {
		return errNoFilename
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errNotExists
		}
		c.Logger().Error("get file info: ", err)
		return errInternal
	}

	if info.IsDir() {
		if err = os.RemoveAll(path); err == nil {
			return c.JSON(200, msg{"Directory deleted successfully"})
		}
	} else {
		if err = os.Remove(path); err == nil {
			return c.JSON(200, msg{"File deleted successfully"})
		}
	}

	c.Logger().Errorf("delete %s: %s", path, err)
	return errInternal
}

func handleMoveFile(c echo.Context) error {
	var act struct{ From, To string }
	if err := c.Bind(&act); err != nil {
		c.Logger().Error("bind act value: ", err)
		return errInternal
	}

	err := os.Rename(filepath.Join(diskRoot, act.From), filepath.Join(diskRoot, act.To))
	if err != nil {
		switch {
		case os.IsExist(err):
			return errExists
		case os.IsNotExist(err):
			return errNotExists
		default:
			c.Logger().Error("rename file: ", err)
			return errInternal
		}
	}

	return c.JSON(200, msg{"File or directory moved successfully"})
}
