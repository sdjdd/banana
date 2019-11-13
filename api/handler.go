package main

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync/atomic"

	"github.com/labstack/echo/v4"
)

func handleGetInfo(c echo.Context) error {
	var userInfo userInfo
	user := c.Get("user")
	if user == nil {
		userInfo.Upload = conf.Anonymous.Upload
		userInfo.Delete = conf.Anonymous.Delete
	} else {
		u := user.(userConf)
		userInfo.Name = c.Get("username").(string)
		userInfo.Upload = u.Upload
		userInfo.Delete = u.Delete
	}

	return c.JSON(200, info{
		Version:  version,
		Capacity: conf.CapBytes,
		Used:     atomic.LoadInt64(&used),
		User:     userInfo,
	})
}

func handleGetFile(c echo.Context) error {
	req := c.Request()
	requestPath := path.Join(conf.Root, req.URL.Path)
	info, err := os.Stat(requestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errNotExists
		}
		c.Logger().Error("get file info:", err)
		return errInternal
	}
	if !info.IsDir() {
		if c.QueryParam("seek") != "" {
			return c.File(requestPath)
		}
		return c.Attachment(requestPath, path.Base(requestPath))
	}

	childrenInfo, err := ioutil.ReadDir(requestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return c.NoContent(404)
		}
		c.Logger().Errorf("get file info: %s", err)
		return errInternal
	}

	result := make([]fileInfo, len(childrenInfo))
	for i, info := range childrenInfo {
		result[i] = fileInfo{
			Name:  info.Name(),
			IsDir: info.IsDir(),
			Size:  info.Size(),
		}
	}
	return c.JSON(200, result)
}

func handlePostFile(c echo.Context) error {
	req := c.Request()
	requestPath := path.Join(conf.Root, req.URL.Path)

	if path.Base(requestPath) == conf.Root {
		return errCreateRoot
	}

	if dirInfo, err := os.Stat(path.Dir(requestPath)); err != nil {
		if os.IsNotExist(err) {
			return errNotExists
		}
		c.Logger().Error("get file info:", err)
		return err
	} else if !dirInfo.IsDir() {
		return errIsNotDir
	}

	if c.QueryParam("type") == "dir" {
		if err := os.Mkdir(requestPath, 0755); err != nil {
			c.Logger().Error("create directory:", err)
			return errInternal
		}
		return c.JSON(200, msg{"directory created successfully"})
	}

	var freeSize int64
	if conf.CapBytes > 0 {
		freeSize = conf.CapBytes - atomic.LoadInt64(&used)
		if freeSize <= 0 {
			return errInsufficient
		}
	}

	f, err := os.OpenFile(requestPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		c.Logger().Error("create file:", err)
		return errInternal
	}

	defer func() {
		f.Close()
		req.Body.Close()
	}()

	if conf.CapBytes > 0 {
		var bufSize int64 = 1024 * 1024 * 4
		var copied int64
		for {
			written, err := io.CopyN(f, req.Body, bufSize)
			if copied += written; copied > freeSize {
				io.Copy(ioutil.Discard, req.Body)
				f.Close()
				os.Remove(requestPath)
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
		atomic.AddInt64(&used, copied)
	} else {
		if _, err := io.Copy(f, req.Body); err != nil {
			c.Logger().Error("copy file:", err)
			return errInternal
		}
	}

	return c.JSON(200, msg{"file uploaded successfully"})
}

func handleDelFile(c echo.Context) error {
	user := c.Get("user")
	if user == nil {
		if !conf.Anonymous.Delete {
			return errForbidden
		}
	} else if !user.(userConf).Delete {
		return errForbidden
	}

	req := c.Request()
	path := filepath.Join(conf.Root, req.URL.Path)

	if filepath.Base(path) == conf.Root {
		return errDeleteRoot
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
