package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const ColonRune = rune(58)

type Command struct {
	Name           string
	Command        string
	CleanupCommand string `yaml:"cleanupCommand"`
	Dir            string
}

type Defs struct {
	Commands []Command
}

func main() {
	defsFile, err := os.ReadFile("defs.yaml")
	if err != nil {
		log.Fatalf("Could not read defs file: %v\n", err)
	}

	defs := Defs{}
	err = yaml.Unmarshal(defsFile, &defs)
	if err != nil {
		log.Fatalf("Could not parse defs file: %v\n", err)
	}

	ui := NewUI()

	commandRunners := make([]*CommandRunner, 0)

	for _, command := range defs.Commands {
		textView := ui.AddLogPage(command.Name)
		commandRunner := NewCommandRunner(
			command.Name,
			command.Command,
			command.CleanupCommand,
			command.Dir,
			textView,
		)
		commandRunner.Run()
		commandRunners = append(commandRunners, commandRunner)
	}

	if err := ui.Run(); err != nil {
		panic(err)
	}

	for _, commandRunner := range commandRunners {
		commandRunner.Stop()
	}
}
