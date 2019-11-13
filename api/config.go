package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type config struct {
	Listen    string
	Root      string
	Capacity  string
	CapBytes  int64
	Anonymous anonymousConf
	Users     map[string]userConf
}

type anonymousConf struct {
	Enable bool
	Upload bool
	Delete bool
}

type userConf struct {
	Password string
	Upload   bool
	Delete   bool
}

var conf config

func init() {
	confile, err := os.Open("banana.yml")
	if err != nil {
		log.Fatalln("open config file:", err)
	}

	err = yaml.NewDecoder(confile).Decode(&conf)
	if err != nil {
		log.Fatalln("parse config file:", err)
	}

	if err := conf.parseCapacity(); err != nil {
		log.Fatalln("parse config file:", err)
	}

	conf.fillDefaults()
	fmt.Printf("%+v\n", conf)
}

func (c *config) fillDefaults() {
	switch {
	case c.Listen == "":
		c.Listen = "0.0.0.0:8080"
	case c.Capacity == "":
		c.Capacity = "unlimited"
	}
}

func (c *config) parseCapacity() error {
	c.Capacity = strings.ToLower(c.Capacity)
	c.CapBytes = 0

	if c.Capacity == "unlimited" {
		return nil
	}

	re := regexp.MustCompile(`^\d+(\.\d+)?[mg]$`)
	if !re.MatchString(c.Capacity) {
		return fmt.Errorf("invalid capacity %q", c.Capacity)
	}

	switch c.Capacity[len(c.Capacity)-1] {
	case 'm':
		c.CapBytes = 1024 * 1024
	case 'g':
		c.CapBytes = 1024 * 1024 * 1024
	}

	n, _ := strconv.ParseFloat(c.Capacity[:len(c.Capacity)-1], 64)
	c.CapBytes = int64(n * float64(c.CapBytes))

	return nil
}
