package cmd

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/vyynl/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, *pokecache.Cache) error
}

// Setting up master command-list
func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a list of valid commands and their functions",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Shows the next 20 location areas to explore",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Shows the last 20 location areas to explore",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore <specific-location-name>",
			description: "Provides a list of potential Pokemon encounters for the specified location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon-name>",
			description: "Attempts to catch the named pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <caught-pokemon-name>",
			description: "Provides the Pokedex entry information for the listed pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all caught pokemon in your Pokedex",
			callback:    commandPokedex,
		},
	}
}

// Declaration of callback functions in the same order as above list
func commandHelp(cfg *Config, cache *pokecache.Cache) error {
	fmt.Printf("\nWelcome to the Pokedex!\nCommand List:\n\n")
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}

func commandExit(cfg *Config, cache *pokecache.Cache) error {
	os.Exit(0)
	return nil
}

func commandMap(cfg *Config, cache *pokecache.Cache) error {
	locationsResp, err := cfg.PokeapiClient.GetListLocations(cfg.NextLocationsURL, cache)
	if err != nil {
		return err
	}

	cfg.NextLocationsURL = locationsResp.Next
	cfg.PrevLocationsURL = locationsResp.Previous

	for _, loc := range locationsResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(cfg *Config, cache *pokecache.Cache) error {
	if cfg.PrevLocationsURL == nil {
		return errors.New("already on first page")
	}

	locationsResp, err := cfg.PokeapiClient.GetListLocations(cfg.PrevLocationsURL, cache)
	if err != nil {
		return err
	}

	cfg.NextLocationsURL = locationsResp.Next
	cfg.PrevLocationsURL = locationsResp.Previous

	for _, loc := range locationsResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandExplore(cfg *Config, cache *pokecache.Cache) error {
	fmt.Printf("Exploring %v\n", cfg.CommandFilter)

	data, err := cfg.PokeapiClient.GetSpecificLocationData(&cfg.CommandFilter, cache)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range data.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *Config, cache *pokecache.Cache) error {
	pokemon, err := cfg.PokeapiClient.GetPokemonData(cfg.CommandFilter, cache)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a ball at %s", pokemon.Name)
	for i := 0; i < 3; i++ {
		time.Sleep(750 * time.Millisecond)
		fmt.Printf(".")
	}
	fmt.Println()

	time.Sleep(1500 * time.Millisecond)
	if pokemon.BaseExperience > rand.IntN(390) {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	} else {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		cfg.Pokedex[pokemon.Name] = pokemon
	}
	return nil
}

func commandInspect(cfg *Config, cache *pokecache.Cache) error {
	if data, exists := cfg.Pokedex[cfg.CommandFilter]; exists {
		fmt.Printf("Name: %s\n", data.Name)
		fmt.Printf("Height: %v\n", data.Height)
		fmt.Printf("Weight: %v\n", data.Weight)
		fmt.Printf("Stats:\n")
		for _, value := range data.Stats {
			fmt.Printf("  - %v: %v\n", value.Stat.Name, value.BaseStat)
		}

		fmt.Printf("Types:")
		fmt.Println()
		for _, value := range data.Types {
			fmt.Printf("  - %v\n", value.Type.Name)
		}
	} else {
		return errors.New("you have not caught that pokemon")
	}
	return nil
}

func commandPokedex(cfg *Config, cache *pokecache.Cache) error {
	fmt.Println("Your Pokedex:")
	for pokemon := range cfg.Pokedex {
		time.Sleep(25 * time.Millisecond)
		fmt.Printf("  - %s\n", pokemon)
	}
	return nil
}
