package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/go-cmd/cmd"
	"github.com/rivo/tview"
)

type CommandRunner struct {
	command        string
	cleanupCommand string
	dir            string
	name           string
	cmd            *cmd.Cmd
	// Channel that is closed when command output is fully read
	readEndChan chan struct{}
	// Channel that yields value when command ends
	endChan  <-chan cmd.Status
	textView *tview.TextView
	exited   bool
}

func NewCommandRunner(name, command, cleanupCommand, dir string, textView *tview.TextView) *CommandRunner {
	cr := new(CommandRunner)
	cr.name = name
	cr.command = command
	cr.cleanupCommand = cleanupCommand
	cr.dir = dir
	cr.readEndChan = make(chan struct{})
	cr.textView = textView

	cr.textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlR {
			cr.Run()
		}

		return event
	})

	return cr
}

func (cr *CommandRunner) Run() {
	cr.exited = false
	cr.readEndChan = make(chan struct{})
	cr.buildCommand()
	cr.endChan = cr.cmd.Start()

	go func() {
		defer close(cr.readEndChan)
		cr.followCommand()
	}()

	go func() {
		<-cr.readEndChan
		cr.exited = true
	}()
}

func (cr *CommandRunner) Stop() error {
	if err := cr.cmd.Stop(); err != nil {
		return err
	}

	// Wait for command to get fully read
	<-cr.readEndChan

	if len(cr.cleanupCommand) > 0 {
		cr.runCleanupCommand()
	}

	return nil
}

func (cr *CommandRunner) GetName() string {
	return cr.name
}

func (cr *CommandRunner) buildCommand() {
	commandParts := strings.Split(cr.command, " ")
	command := cmd.NewCmdOptions(cmd.Options{
		Buffered:  false,
		Streaming: true,
	}, commandParts[0], commandParts[1:]...)
	command.Dir = cr.dir
	cr.cmd = command
}

func (cr *CommandRunner) followCommand() {
	for cr.cmd.Stdout != nil || cr.cmd.Stderr != nil {
		select {
		case line, open := <-cr.cmd.Stdout:
			if !open {
				cr.cmd.Stdout = nil
				continue
			}
			fmt.Fprintln(cr.textView, line)
		case line, open := <-cr.cmd.Stderr:
			if !open {
				cr.cmd.Stderr = nil
				continue
			}
			fmt.Fprintln(cr.textView, line)
		}
	}
}

func (cr *CommandRunner) runCleanupCommand() {
	commandParts := strings.Split(cr.cleanupCommand, " ")
	cmd := cmd.NewCmd(commandParts[0], commandParts[1:]...)
	cmd.Dir = cr.dir
	status := <-cmd.Start()

	for _, line := range status.Stdout {
		fmt.Fprintln(os.Stdout, line)
	}
	for _, line := range status.Stderr {
		fmt.Fprintln(os.Stderr, line)
	}
}
