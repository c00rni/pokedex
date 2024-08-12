package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/c00rni/pokedex/internal/api"
	"github.com/c00rni/pokedex/internal/pokecache"
	"math/rand"
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

type pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height                 int    `json:"height"`
	HeldItems              []any  `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []any  `json:"past_abilities"`
	PastTypes     []any  `json:"past_types"`
	Species       struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       any    `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  any    `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      any    `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale any    `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       any    `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       any    `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  any    `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

type config struct {
	Next     string
	Previous string
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
	pokedex := map[string]pokemon{}

	commandCatch := func(opts ...string) error {
		if len(opts) < 1 {
			return errors.New("The catch command need a pokemon name as argument")
		}
		_, ok := pokedex[opts[0]]
		if ok {
			fmt.Println("Pokemon already captured.")
			return nil
		}
		target := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v", opts[0])
		responseBytes, err := api.GetLocations(target)
		if err != nil {
			return err
		}
		byteBodyResponse := responseBytes
		pokemonDetails := pokemon{}
		err1 := json.Unmarshal(byteBodyResponse, &pokemonDetails)
		if err1 != nil {
			return err1
		}

		fmt.Println(fmt.Sprintf("Throwing a Pokeball at %v...", opts[0]))
		if float64(pokemonDetails.BaseExperience)*rand.NormFloat64() < 10 {
			pokedex[opts[0]] = pokemonDetails
			fmt.Println(fmt.Sprintf("%v was caught!", opts[0]))
		} else {
			fmt.Println(fmt.Sprintf("%v escaped!", opts[0]))
		}
		return nil
	}

	commandInspect := func(opts ...string) error {
		if len(opts) < 1 {
			return errors.New("The catch command need a pokemon name as argument")
		}
		pokemonDetails, ok := pokedex[opts[0]]
		if !ok {
			fmt.Println("you have not caught that pokemon")
			return nil
		}
		detailsString := fmt.Sprintf("Name: %v\nHeight: %v\nWeight: %v\nStats:", pokemonDetails.Forms[0].Name, pokemonDetails.Height, pokemonDetails.Weight)

		fmt.Println(detailsString)
		for _, stats := range pokemonDetails.Stats {
			line := fmt.Sprintf(" - %v: %v", stats.Stat.Name, stats.BaseStat)
			fmt.Println(line)
		}
		fmt.Println("Types:")
		for _, types := range pokemonDetails.Types {
			line := fmt.Sprintf(" - %v", types.Type.Name)
			fmt.Println(line)
		}
		return nil
	}

	commandPokedex := func(opts ...string) error {
		if len(opts) < 0 {
			return errors.New("The pokedex command dont expect arguments.")
		}

		fmt.Println("Your Pokedex:")
		for _, pokemon := range pokedex {
			fmt.Println(" -", pokemon.Forms[0].Name)
		}
		return nil
	}

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

	commands := make(map[string]cliCommand)

	commandHelp := func(opts ...string) error {
		fmt.Println("Welcome to the Pokedex!\nUsage:")
		for _, command := range commands {
			fmt.Println(fmt.Sprintf("%v: %v", command.name, command.description))
		}
		return nil
	}

	scanner := bufio.NewScanner(os.Stdin)

	commands = map[string]cliCommand{
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
			description: "Discover new areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Diplay previous areas",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "List pokemons in an area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to capture a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Print stats about a pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Print all the captured pokemon names",
			callback:    commandPokedex,
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
