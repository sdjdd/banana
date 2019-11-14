package main

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v2"
)

type config struct {
	Listen string
	Root   string
	Size   int64
}

type user struct {
	Expire    time.Time
	Password  string
	Privilege map[string]struct{}
}

type userList map[string]user

var conf config

var users = make(userList)

func compare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func (u userList) auth(username, password string, c echo.Context) (bool, error) {
	var user user
	var ok bool
	if username == "" || username == "anonymous" {
		if user, ok = u["anonymous"]; !ok {
			return false, nil
		}
	} else {
		if user, ok = u[username]; !ok || !compare(user.Password, password) {
			return false, nil
		}
	}

	if !user.Expire.IsZero() && time.Now().After(user.Expire) {
		return false, nil
	}

	c.Set("user", user)
	return true, nil
}

func (u user) can(action string) bool {
	_, ok := u.Privilege[action]
	return ok
}

func (c *config) read(name string) (err error) {
	search := []string{name, filepath.Join(binDir, name)}
	var f *os.File
	for _, path := range search {
		f, err = os.Open(path)
		if err == nil {
			break
		}
	}
	if err != nil {
		return
	}
	defer f.Close()

	var rawConf struct {
		Listen string
		Root   string
		Size   string
		Users  map[string]struct {
			Expire    string
			Password  string
			Privilege []string
		}
	}

	err = yaml.NewDecoder(f).Decode(&rawConf)
	if err != nil {
		return err
	} else if rawConf.Root == "" {
		return errors.New("root must be provided")
	}

	if rawConf.Listen == "" {
		c.Listen = "0.0.0.0:8080"
	} else {
		c.Listen = rawConf.Listen
	}

	if c.Root, err = filepath.Abs(rawConf.Root); err != nil {
		return
	}
	if c.Size, err = parseSize(rawConf.Size); err != nil {
		return fmt.Errorf("parse size %q: %s", rawConf.Size, err)
	}

	for name, info := range rawConf.Users {
		privilege := make(map[string]struct{})
		for _, pri := range info.Privilege {
			privilege[pri] = struct{}{}
		}

		var ex time.Time
		if info.Expire != "" {
			ex, err = time.ParseInLocation("2006-01-02 15:04:05",
				info.Expire, time.Local)
			if err != nil {
				return
			}
		}

		users[name] = user{
			Expire:    ex,
			Password:  info.Password,
			Privilege: privilege,
		}
	}

	return nil
}

func parseSize(str string) (size int64, err error) {
	if str == "" {
		return 0, nil
	}

	var pointed bool
	buf := new(strings.Builder)
	for i, ch := range str {
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		} else {
			switch ch {
			case '.':
				if pointed {
					err = fmt.Errorf("unexpected '.' at %d", i+1)
					return
				}
				pointed = true
				buf.WriteRune(ch)
			case 'k', 'm', 'g':
				fallthrough
			case 'K', 'M', 'G':
				if i != len(str)-1 {
					err = fmt.Errorf("unexpected %q at %d", str[i+1], i+2)
					return
				}
			default:
				err = fmt.Errorf("unexpected %q at %d", str[i], i+1)
				return
			}
		}
	}

	n, _ := strconv.ParseFloat(buf.String(), 64)
	switch str[len(str)-1] {
	case 'k', 'K':
		size = 1024
	case 'm', 'M':
		size = 1024 * 1024
	case 'g', 'G':
		size = 1024 * 1024 * 1024
	default:
		size = 1
	}

	size = int64(float64(size) * n)
	return
}
