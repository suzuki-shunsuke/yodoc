package main

import (
	"os"

	"github.com/suzuki-shunsuke/yodoc/pkg/cli"
)

var version = ""

func main() {
	if code := cli.Run(version); code != 0 {
		os.Exit(code)
	}
}
