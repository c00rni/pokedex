package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/c00rni/pokedex/internal/api"
	"github.com/c00rni/pokedex/internal/pokecache"
	"os"
	"strings"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(...string) error
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

type area struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type config struct {
	Next     string
	Previous string
}

func commandHelp(opts ...string) error {
	_, err := fmt.Print(`
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
map: List 20 more locations
mapb: List 20 last locations
explore: List pokemon foun from an area
`)
	return err
}

func commandExit(_ ...string) error {
	os.Exit(0)
	return nil
}

func printLocations(response response) {
	for _, result := range response.Results {
		fmt.Println(result.Name)
	}
}

func printPokemons(area area) {
	fmt.Println("Exploring pastoria-city-area...")
	fmt.Println("Found Pokemon:")
	for _, data := range area.PokemonEncounters {
		fmt.Println(" -", data.Pokemon.Name)
	}
}

func main() {
	config := config{
		Next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		Previous: "",
	}

	interval := time.Minute
	catch := pokecache.NewCache(interval)

	commandExplore := func(opts ...string) error {
		var byteBodyResponse []byte
		if len(opts) < 1 {
			return errors.New("The explore command needs one area name")
		}
		target := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v/", opts[0])
		cacheBytes, ok := catch.Get(target)
		if ok {
			byteBodyResponse = cacheBytes
		} else {
			responseBytes, err := api.GetLocations(target)
			if err != nil {
				return err
			}
			byteBodyResponse = responseBytes
			catch.Add(target, byteBodyResponse)
		}
		areaDetails := area{}
		err1 := json.Unmarshal(byteBodyResponse, &areaDetails)
		if err1 != nil {
			return err1
		}

		printPokemons(areaDetails)
		return nil
	}

	commandMap := func(_ ...string) error {
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

	commandMapB := func(_ ...string) error {
		if config.Previous == "" {
			config.Next = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
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
		"explore": {
			name:        "explore",
			description: "List pokemons in an area",
			callback:    commandExplore,
		},
	}

	fmt.Print("pokedex > ")
	for scanner.Scan() {
		inputs := strings.Split(scanner.Text(), " ")
		if cmd, ok := commands[inputs[0]]; ok {
			cmd.callback(inputs[1:]...)
		}
		fmt.Print("pokedex > ")
	}
}
