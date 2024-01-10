package cmd

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

type Command struct {
	Http       bool
	Version    bool
	ForceIpv4  bool
	ServerPort int
	Loglevel   string
}

var (
	c    *Command
	once sync.Once
)

func Get() *Command {
	once.Do(func() {
		c = new(Command)
		flag.IntVar(&c.ServerPort, "ap", 10000, "api server port")
		flag.StringVar(&c.Loglevel, "ll", "info", "set loglevel")
		flag.BoolVar(&c.Http, "http", false, "http mode")
		flag.BoolVar(&c.ForceIpv4, "ipv4", false, "force ipv4")
		flag.BoolVar(&c.Version, "v", false, "version")
		flag.Parse()
		if err := c.check(); err != nil {
			fmt.Println("command check:", err)
			os.Exit(0)
		}
	})
	return c
}

func (c *Command) Debug() bool {
	return c.Loglevel == "debug"
}

func (c *Command) check() error {
	if c.ServerPort < 1 || c.ServerPort > 65535 {
		return fmt.Errorf("port set error, range is [1, 65535], current port is %d", c.ServerPort)
	}
	if ss := []string{"debug", "info", "warn", "error"}; !c.include(ss, c.Loglevel) {
		return fmt.Errorf("loglevel should from %v", ss)
	}
	return nil
}

func (c *Command) include(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
