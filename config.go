package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	reInclude = regexp.MustCompile(
		`(.*)?!include\s+(.*)`,
	)
)

type Config struct {
	Actions map[string]Action `yaml:"actions"`
	Rules   []Rule            `yaml:"rules"`
}

func getConfig(path string) (*Config, error) {
	yamlData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	yamlData, err = invokeIncludes(yamlData, filepath.Dir(path))
	if err != nil {
		return nil, err
	}

	//config := map[string]interface{}{}
	config := &Config{}

	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		return nil, fmt.Errorf("can't decode yaml: %s", err)
	}

	//read name from key of action directive
	for actionName, action := range config.Actions {
		action.Name = actionName
		config.Actions[actionName] = action
	}

	return config, nil
}

func invokeIncludes(yamlData []byte, baseDir string) ([]byte, error) {
	yamlString := string(yamlData)

	matches := reInclude.FindAllStringSubmatch(yamlString, -1)
	if len(matches) == 0 {
		return yamlData, nil
	}

	for _, match := range matches {
		var (
			rawLine   = match[0]
			rawIndent = match[1]
			pattern   = match[2]
		)

		if !strings.HasPrefix(pattern, "/") {
			pattern = filepath.Join(baseDir, pattern)
		}

		files, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf(
				"can't glob by pattern '%s': %s",
				pattern, err,
			)
		}

		var including string
		for _, file := range files {
			fileData, err := ioutil.ReadFile(file)
			if err != nil {
				return nil, fmt.Errorf(
					"can't include by pattern '%s': %s", pattern, err,
				)
			}

			including = string(fileData) + "\n"
		}

		if len(rawIndent) > 0 {
			level := 0
			nextLevel := false
			if string(rawIndent[0]) == " " {
				trimmed := strings.TrimLeft(rawIndent, " ")
				spaces := len(rawIndent) - len(trimmed)
				level = spaces / 4

				if trimmed != "" {
					level = level + 1
					nextLevel = true
				}
			}

			lines := strings.Split(including, "\n")
			for i, line := range lines {
				if !nextLevel && i == 0 {
					continue
				}

				lines[i] = strings.Repeat(" ", level*4) + line
			}

			including = strings.Join(lines, "\n")
			if nextLevel {
				including = "\n" + including
			}

			including = rawIndent + including
		}

		yamlString = strings.Replace(yamlString, rawLine, including, -1)
	}

	return []byte(yamlString), nil
}
