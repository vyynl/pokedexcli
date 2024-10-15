package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/vyynl/pokedexcli/internal/pokeapi"
	"github.com/vyynl/pokedexcli/internal/pokecache"
)

type Config struct {
	PokeapiClient    pokeapi.Client
	NextLocationsURL *string
	PrevLocationsURL *string
	CommandFilter    string
	Pokedex          map[string]pokeapi.Pokemon
}

func StartRepl(cfg *Config, cache *pokecache.Cache) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Pokedex > ")
		scanner.Scan()

		inputCommand := cleanInput(scanner.Text())
		commandName := inputCommand[0]

		if len(commandName) == 0 {
			continue
		}

		if len(inputCommand) > 1 {
			cfg.CommandFilter = inputCommand[1]
		} else {
			cfg.CommandFilter = ""
		}

		command, exists := getCommands()[commandName]
		if exists {
			err := command.callback(cfg, cache)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Command not found")
			continue
		}
	}
}

func cleanInput(input string) []string {
	lowerInput := strings.ToLower(input)
	return strings.Fields(lowerInput)
}
