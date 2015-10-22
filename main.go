package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
)

var (
	version = "1.0"
)

const (
	usage = `Guntalina

Usage:
	guntalina [options] -s <source> -c <config>

Options:
	-s <source>    Source file
	-c <config     Config file
	-r --dry-run   Dry-run mode, in this mode commands will be not really executed.
	-f --force     Do not stop if any command has been failed.
`
)

func main() {
	args, err := docopt.Parse(usage, nil, true, version, true, true)
	if err != nil {
		panic(err)
	}

	var (
		modificationsPath = args["-s"].(string)
		configPath        = args["-c"].(string)
		dryRun            = args["--dry-run"].(bool)
		force             = args["--force"].(bool)
	)

	modificationsData, err := ioutil.ReadFile(modificationsPath)
	if err != nil {
		log.Fatal(err)
	}

	modifications := strings.Split(
		strings.TrimSuffix(string(modificationsData), "\n"), "\n",
	)

	config, err := getConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	actions, rules, err := config.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	workflow := []string{}
	for _, modification := range modifications {
		rule, err := rules.GetRule(modification)
		if err != nil {
			log.Fatal(err)
		}

		if rule == nil {
			continue
		}

		workflow = append(workflow, rule.Workflow...)

		rules.SetProcessed(rule)
	}

	if len(workflow) == 0 {
		log.Printf("nothing to do")
		return
	}

	// prevent double execution
	workflow = uniqueWorkflow(workflow)

	log.Println("following actions will be executed:")

	commands := []string{}
	for _, actionName := range workflow {
		action, ok := actions[actionName]
		if !ok {
			log.Fatal(
				"[BUG] can't find action '%s' in action array, "+
					"possible validation error",
				actionName,
			)
		}

		log.Println(actionName)
		for _, command := range action.Commands {
			log.Printf("    %s\n", command)
		}
		log.Println()

		commands = append(commands, action.Commands...)
	}

	log.Println("following commands will be executed:")
	for _, command := range commands {
		log.Println(command)
	}
	log.Println()

	if dryRun {
		return
	}

	for _, command := range commands {
		log.Printf("executing: %s\n", command)

		output, err := execute(command)
		if len(output) != 0 {
			log.Println(output)
		}

		if err != nil {
			log.Println(err.Error())
			log.Println() // add empty line for pretty output
			if !force {
				os.Exit(1)
			}
		}
	}
}

func uniqueWorkflow(workflow []string) []string {
	unique := []string{}
	for _, action := range workflow {
		found := false
		for _, item := range unique {
			if item == action {
				found = true
				break
			}
		}

		if found {
			continue
		}

		unique = append(unique, action)
	}

	return unique
}
