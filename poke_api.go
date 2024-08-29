package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "https://pokeapi.co/api/v2"
)

func newClient(cacheInterval time.Duration) client {
	return client{
		cache:  NewCache(cacheInterval),
		client: http.Client{},
	}
}

func (c *client) ListLocations(pageURL *string) (RespShallowLocations, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	if val, ok := c.cache.Get(url); ok {
		locationsResp := RespShallowLocations{}
		err := json.Unmarshal(val, &locationsResp)
		if err != nil {
			return RespShallowLocations{}, err
		}

		return locationsResp, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RespShallowLocations{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return RespShallowLocations{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RespShallowLocations{}, err
	}

	locationsResp := RespShallowLocations{}
	err = json.Unmarshal(dat, &locationsResp)
	if err != nil {
		return RespShallowLocations{}, err
	}

	return locationsResp, nil
}

func (c *client) GetLocation(locationName string) (Location, error) {
	url := baseURL + "/location-area/" + locationName

	if val, ok := c.cache.Get(url); ok {
		response := Location{}
		err := json.Unmarshal(val, &response)
		if err != nil {
			return Location{}, err
		}
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Location{}, err
	}

	response, err := c.client.Do(request)
	if err != nil {
		return Location{}, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return Location{}, err
	}

	locationResponse := Location{}
	err = json.Unmarshal(data, &locationResponse)
	if err != nil {
		return Location{}, err
	}

	c.cache.Add(url, data)

	return locationResponse, nil
}

func (c *client) GetPokemonInfo(pokemonName string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + pokemonName

	if val, ok := c.cache.Get(url); ok {
		response := Pokemon{}
		err := json.Unmarshal(val, &response)
		if err != nil {
			return Pokemon{}, err
		}
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Pokemon{}, err
	}

	response, err := c.client.Do(request)
	if err != nil {
		return Pokemon{}, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return Pokemon{}, err
	}

	pokemonData := Pokemon{}
	err = json.Unmarshal(data, &pokemonData)
	if err != nil {
		return Pokemon{}, err
	}

	c.cache.Add(url, data)

	return pokemonData, nil
}
