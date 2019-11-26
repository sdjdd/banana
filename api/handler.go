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

type msg struct {
	Message string `json:"message"`
}

type fileInfo struct {
	Name    string    `json:"name"`
	IsDir   bool      `json:"isDir"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
}

const bufSize int64 = 1024 * 1024 * 4

func handleGetInfo(c echo.Context) error {
	return c.JSON(200, struct {
		Version string `json:"version"`
		Size    int64  `json:"size"`
		Used    int64  `json:"used"`
	}{
		Version: version,
		Size:    conf.Size,
		Used:    atomic.LoadInt64(&used),
	})
}

func handleGetFile(c echo.Context) error {
	path := c.Get("path").(string)
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errNotExists
		}
		c.Logger().Error("get file info:", err)
		return errInternal
	}

	if !info.IsDir() {
		user := c.Get("user").(*user)
		if !user.Privilege.Download {
			return errForbidden
		}
		return c.Attachment(path, filepath.Base(path))
	}

	childrenInfo, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errNotExists
		}
		c.Logger().Errorf("get file info: %s", err)
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
	user := c.Get("user").(*user)
	if !user.Privilege.Upload {
		return errForbidden
	}

	path := c.Get("path").(string)
	if path == conf.Root {
		return errNoFilename
	}

	if dirInfo, err := os.Stat(filepath.Dir(path)); err != nil {
		if os.IsNotExist(err) {
			return errBaseNotExists
		}
		c.Logger().Error("get base info:", err)
		return err
	} else if !dirInfo.IsDir() {
		return errIsNotDir
	}

	if c.QueryParam("type") == "dir" {
		if err := os.Mkdir(path, 0755); err != nil {
			if os.IsExist(err) {
				return errExists
			}
			c.Logger().Error("create directory:", err)
			return errInternal
		}
		return c.JSON(200, msg{"directory created successfully"})
	}

	var freeSize int64
	if conf.Size > 0 {
		freeSize = conf.Size - atomic.LoadInt64(&used)
		if freeSize <= 0 {
			return errInsufficient
		}
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		c.Logger().Error("create file:", err)
		return errInternal
	}

	body := c.Request().Body
	defer f.Close()

	var copied int64
	if conf.Size > 0 {
		for {
			written, err := io.CopyN(f, body, bufSize)
			if copied += written; copied > freeSize {
				f.Close()
				os.Remove(path)
				return errInsufficient
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				c.Logger().Error("copy file:", err)
				return errInternal
			}
		}
	} else {
		if copied, err = io.Copy(f, body); err != nil {
			c.Logger().Error("copy file:", err)
			return errInternal
		}
	}
	atomic.AddInt64(&used, copied)

	return c.JSON(200, msg{"file uploaded successfully"})
}

func handleDelFile(c echo.Context) error {
	user := c.Get("user").(*user)
	if !user.Privilege.Delete {
		return errForbidden
	}

	path := c.Get("path").(string)
	if path == conf.Root {
		return errNoFilename
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errNotExists
		}
		c.Logger().Errorf("get file info: %s", err)
		return errInternal
	}

	if info.IsDir() {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}
	if err != nil {
		c.Logger().Errorf("delete file or directory: %s", err)
		return errInternal
	}

	return c.JSON(200, msg{"file deleted successfully"})
}

func handleMoveFile(c echo.Context) error {
	user := c.Get("user").(*user)
	if !user.Privilege.Upload || !user.Privilege.Delete {
		return errForbidden
	}

	var act struct{ From, To string }
	if err := c.Bind(&act); err != nil {
		c.Logger().Errorf("bind act value: %s", err)
		return errInternal
	}

	act.From = filepath.Join(conf.Root, act.From)
	act.To = filepath.Join(conf.Root, act.To)

	if err := os.Rename(act.From, act.To); err != nil {
		switch {
		case os.IsExist(err):
			return errExists
		case os.IsNotExist(err):
			return errNotExists
		default:
			c.Logger().Errorf("rename file: %s", err)
			return errInternal
		}
	}

	return c.JSON(200, msg{"file or directory moved successfully"})
}
