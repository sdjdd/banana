package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type config struct {
	Listen string
	Root   string
	Size   int64
}

type confUsers map[string]struct {
	Expire    string
	Password  string
	Privilege []string
}

func (c *config) reset() {
	c.Listen = ""
	c.Root = ""
	c.Size = 0
}

func (c *config) read(name string) (err error) {
	c.reset()

	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	var rawConf struct {
		Listen string
		Root   string
		Size   string
		Users  confUsers
	}

	if err = yaml.NewDecoder(f).Decode(&rawConf); err != nil {
		return
	} else if rawConf.Root == "" {
		return errors.New("root must be provided")
	}

	if c.Listen = rawConf.Listen; c.Listen == "" {
		c.Listen = "0.0.0.0:8080"
	}
	if c.Root, err = filepath.Abs(rawConf.Root); err != nil {
		return
	}
	if c.Size, err = parseSize(rawConf.Size); err != nil {
		return fmt.Errorf("parse size %q: %s", rawConf.Size, err)
	}

	return users.Set(rawConf.Users)
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
