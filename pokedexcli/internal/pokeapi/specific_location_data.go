package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/vyynl/pokedexcli/internal/pokecache"
)

// Capture Location Data
func (c *Client) GetSpecificLocationData(specificLocation *string, cache *pokecache.Cache) (RespSpecificLocation, error) {
	if len(*specificLocation) == 0 {
		return RespSpecificLocation{}, errors.New("please enter a location to explore")
	}

	url := baseURL + "/location-area/" + *specificLocation
	// 1: Checking cache
	if cachedRes, exists := cache.Get(url); exists {
		var SpecificLocationResp RespSpecificLocation
		err := json.Unmarshal(cachedRes, &SpecificLocationResp)
		if err != nil {
			return RespSpecificLocation{}, err
		}
		return SpecificLocationResp, nil
	}

	// 2: Fetching from network if not cached
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RespSpecificLocation{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return RespSpecificLocation{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return RespSpecificLocation{}, err
	}

	// 3: Adding new data to cache
	cache.Add(url, data)

	// 4: Parse and return the data
	specificLocationResp := RespSpecificLocation{}
	err = json.Unmarshal(data, &specificLocationResp)
	if err != nil {
		return RespSpecificLocation{}, err
	}

	return specificLocationResp, nil
}
