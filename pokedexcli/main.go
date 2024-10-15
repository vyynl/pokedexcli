package main

import (
	"time"

	"github.com/vyynl/pokedexcli/cmd"
	"github.com/vyynl/pokedexcli/internal/pokeapi"
	"github.com/vyynl/pokedexcli/internal/pokecache"
)

func main() {
	pokeClient := pokeapi.NewClient(5 * time.Second)
	cfg := &cmd.Config{
		PokeapiClient: pokeClient,
		Pokedex:       make(map[string]pokeapi.Pokemon),
	}
	cache := pokecache.NewCache(5 * time.Second)

	cmd.StartRepl(cfg, cache)
}
