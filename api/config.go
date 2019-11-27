package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type config struct {
	Listen string
	Root   string
	Size   string
	Users  confUsers
}

type confUsers map[string]struct {
	Expire    string
	Password  string
	Privilege []string
}

func readConfig(name string) (c config, err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	if err = yaml.NewDecoder(f).Decode(&c); err != nil {
		return
	}
	if c.Listen == "" {
		c.Listen = "0.0.0.0:8080"
	}
	if c.Root == "" {
		err = errors.New("root must be provided")
		return
	}

	return
}

func parseSize(str string) (size int64, err error) {
	origin := str
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return 0, nil
	}
	str = strings.ToUpper(str)

	re := regexp.MustCompile(`^\d+(\.\d+)?[KMG]?B?$`)
	if !re.MatchString(str) {
		err = fmt.Errorf("parsing %q: invalid syntax", origin)
		return
	}

	if str[len(str)-1] == 'B' {
		str = str[:len(str)-1]
	}

	switch str[len(str)-1] {
	case 'K':
		size = 1 << 10
		str = str[:len(str)-1]
	case 'M':
		size = 1 << 20
		str = str[:len(str)-1]
	case 'G':
		size = 1 << 30
		str = str[:len(str)-1]
	default:
		size = 1
	}

	num, err := strconv.ParseFloat(str, 64)
	size = int64(float64(size) * num)
	return
}
