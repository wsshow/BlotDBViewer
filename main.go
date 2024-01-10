package main

import (
	"BBoltViewer/cmd"
	"BBoltViewer/g"
	"BBoltViewer/server"
	"BBoltViewer/version"
	"fmt"
)

func main() {
	c := cmd.Get()
	if c.Version {
		fmt.Println(version.Get())
		return
	}
	g.Conf = c
	g.Log = g.NewLog()
	server.Run(c)
}
