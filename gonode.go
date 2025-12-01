package main

import (
	"github.com/VastBlast/gonode/cli"
)

const (
	Name    = "gonode"
	Version = "v1.2.2"
)

func main() {
	c := cli.CLI{}
	c.Run(Name, Version)
}
