package main

type Action struct {
	Name     string   `yaml:"name"`
	Commands []string `yaml:"commands"`
}
