package main

import (
	"os"

	"github.com/ADXenomorph/go-everyqa/cli"
)

func main() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	cli.Run()
}
