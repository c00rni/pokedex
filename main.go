package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/c00rni/pokedex/internal/api"
	"github.com/c00rni/pokedex/internal/pokecache"
	"os"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type response struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type config struct {
	Next     string
	Previous string
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

func printLocations(response response) {
	for _, result := range response.Results {
		fmt.Println(result.Name)
	}
}

func main() {
	config := config{
		Next:     "https://pokeapi.co/api/v2/location/?offset=0&limit=20",
		Previous: "",
	}

	interval := time.Minute
	catch := pokecache.NewCache(interval)

	commandMap := func() error {
		var byteBodyResponse []byte
		cacheBytes, ok := catch.Get(config.Next)
		if ok {
			byteBodyResponse = cacheBytes
		} else {
			responseBytes, err := api.GetLocations(config.Next)
			if err != nil {
				return err
			}
			byteBodyResponse = responseBytes
		}
		catch.Add(config.Next, byteBodyResponse)
		response := response{}
		err1 := json.Unmarshal(byteBodyResponse, &response)
		if err1 != nil {
			return err1
		}
		config.Next = response.Next
		config.Previous = response.Previous

		printLocations(response)
		return nil
	}

	commandMapB := func() error {
		if config.Previous == "" {
			config.Next = "https://pokeapi.co/api/v2/location/?offset=0&limit=20"
			return nil
		}
		var byteBodyResponse []byte
		cacheBytes, ok := catch.Get(config.Previous)
		if ok {
			byteBodyResponse = cacheBytes
		} else {
			responseBytes, err := api.GetLocations(config.Previous)
			if err != nil {
				return err
			}
			byteBodyResponse = responseBytes
		}
		catch.Add(config.Previous, byteBodyResponse)
		response := response{}
		err1 := json.Unmarshal(byteBodyResponse, &response)
		if err1 != nil {
			return err1
		}
		config.Next = response.Next
		config.Previous = response.Previous

		printLocations(response)
		return nil
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
