package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandHelp() error {
	_, err := fmt.Print(`
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex

`)
	return err
}

func commandExit() error {
	os.Exit(0)
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	commands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	fmt.Print("pokedex > ")
	for scanner.Scan() {

		if cmd, ok := commands[scanner.Text()]; ok {
			cmd.callback()
		}
		fmt.Print("pokedex > ")
	}
}
