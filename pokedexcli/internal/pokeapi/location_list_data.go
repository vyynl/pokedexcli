package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/vyynl/pokedexcli/internal/pokecache"
)

// List Locations
func (c *Client) GetListLocations(pageURL *string, cache *pokecache.Cache) (RespShallowLocations, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	// 1: Checking Cache
	if cachedRes, exists := cache.Get(url); exists {
		var locationsResp RespShallowLocations
		err := json.Unmarshal(cachedRes, &locationsResp)
		if err != nil {
			return RespShallowLocations{}, err
		}
		return locationsResp, nil
	}

	// 2: Fetching from network if not cached
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RespShallowLocations{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return RespShallowLocations{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return RespShallowLocations{}, err
	}

	// 3: Adding new data to cache
	cache.Add(url, data)

	// 4: Parse and return the data
	locationsResp := RespShallowLocations{}
	err = json.Unmarshal(data, &locationsResp)
	if err != nil {
		return RespShallowLocations{}, err
	}

	return locationsResp, nil
}
