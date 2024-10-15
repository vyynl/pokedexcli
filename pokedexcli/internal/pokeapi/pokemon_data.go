package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/vyynl/pokedexcli/internal/pokecache"
)

// Collecting the data for the Pokemon that we're trying to catch
func (c *Client) GetPokemonData(pokemonName string, cache *pokecache.Cache) (Pokemon, error) {
	url := baseURL + "/pokemon/" + pokemonName
	if len(pokemonName) == 0 {
		return Pokemon{}, errors.New("please enter a pokemon to catch")
	}

	// 1: Checking Cache
	if cachedRes, exists := cache.Get(url); exists {
		pokemonData := Pokemon{}
		err := json.Unmarshal(cachedRes, &pokemonData)
		if err != nil {
			return Pokemon{}, err
		}
		return pokemonData, nil
	}

	// 2: Fetching from network if not cached
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Pokemon{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}

	// 3: Adding new data to cache
	cache.Add(url, data)

	// 4: Parse and return the data
	pokemonData := Pokemon{}
	err = json.Unmarshal(data, &pokemonData)
	if err != nil {
		return Pokemon{}, err
	}

	return pokemonData, nil
}
