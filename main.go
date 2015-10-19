package main

import (
	"log"

	"github.com/docopt/docopt-go"
)

var (
	version = "1.0"
)

const (
	usage = `Guntalina

Usage:
	./guntalina -s <source> -c <config>

Options:
	-s <source>    Source file
	-c <config     Config file

`
)

func main() {
	args, err := docopt.Parse(usage, nil, true, version, true, true)
	if err != nil {
		panic(err)
	}

	var (
		sourcePath = args["-s"].(string)
		configPath = args["-c"].(string)
	)

	_, err = getConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
}
