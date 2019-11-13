package main

import (
	"crypto/subtle"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync/atomic"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type fileInfo struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
}

type info struct {
	Version  string `json:"version"`
	Capacity int64  `json:"capacity"`
	Used     int64  `json:"used"`
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
	errForbidden    = echo.NewHTTPError(403, httpErr{1, "forbidden"})
	errCreateRoot   = echo.NewHTTPError(400, httpErr{2, "cannot create /"})
	errNotExists    = echo.NewHTTPError(404, httpErr{3, "no such file of directiry"})
	errIsNotDir     = echo.NewHTTPError(400, httpErr{4, "is not directory"})
	errInsufficient = echo.NewHTTPError(400, httpErr{5, "space is insufficient"})
)

var used int64

func main() {
	e := echo.New()
	e.HideBanner = true
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
	if conf.Anonymous.Enable {
		c.Set("username", "")
		return true
	}
	return false
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
	return true, nil
}

func handleGetInfo(c echo.Context) error {
	return c.JSON(200, info{
		Version:  version,
		Capacity: conf.CapBytes,
		Used:     atomic.LoadInt64(&used),
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
		c.Logger().Error("get fileinfo:", err)
		return err
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
		// log.Printf("read dir %q: %s\n", path, err)
		return err
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
		c.Logger().Error("get fileinfo:", err)
		return err
	} else if !dirInfo.IsDir() {
		return errIsNotDir
	}

	if c.QueryParam("type") == "dir" {
		if err := os.Mkdir(requestPath, 0755); err != nil {
			c.Logger().Error("create directory:", err)
			return err
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
		return err
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
				return err
			}
		}
		atomic.AddInt64(&used, copied)
	} else {
		if _, err := io.Copy(f, req.Body); err != nil {
			c.Logger().Error("copy file:", err)
			return err
		}
	}

	return c.JSON(200, msg{"file uploaded successfully"})
}
