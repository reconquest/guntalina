package main

type Rule struct {
	Masks   []string `yaml:"masks"`
	Actions []string `yaml:"actions"`
}
