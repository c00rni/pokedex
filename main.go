package main

import (
	"bufio"
	"fmt"
	"github.com/c00rni/pokedex/internal/api"
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
map: List 20 more locations
mapb: List 20 last locations
`)
	return err
}

func commandExit() error {
	os.Exit(0)
	return nil
}

func main() {
	config := api.Config{
		Next:     "https://pokeapi.co/api/v2/location/?offset=0&limit=20",
		Previous: "",
	}

	commandMap := func() error {
		return api.GetLocations(true, &config)
	}

	commandMapB := func() error {
		return api.GetLocations(false, &config)
	}

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
		"map": {
			name:        "map",
			description: "Explore 20 more locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Explore 20 last locations",
			callback:    commandMapB,
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
