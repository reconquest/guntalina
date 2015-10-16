package main

import (
	"fmt"

	"github.com/docopt/docopt-go"
)

var (
	version = "1.0"
)

const (
	usage = `Guntalina

Usage:
	guntalina -s <path> -r <path>
`
)

func main() {
	args, err := docopt.Parse(usage, nil, true, version, true, true)
	if err != nil {
		panic(err)
	}

	var (
		sourcePath = args["-s"].(string)
		rulesPath  = args["-r"].(string)
	)

	fmt.Printf("XXXXXX main.go:29: sourcePath: %#v\n", sourcePath)
	fmt.Printf("XXXXXX main.go:29: rulesPath: %#v\n", rulesPath)
}
