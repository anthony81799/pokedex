package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/anthony81799/pokedex/config"
)

func GetMapLocations(cfg *config.Config, url string) ([]Location, error) {
	res, err := http.Get(url)
	if err != nil {
		return []Location{}, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return []Location{}, fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return []Location{}, err
	}

	locations := LocationAreaRes{}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return []Location{}, err
	}

	// Safely handle the Next pointer
	if locations.Next != nil {
		cfg.Next = locations.Next
	} else {
		cfg.Next = nil // or handle end of pagination appropriately
	}

	// Previous can also be nil, handle it safely
	if locations.Previous != nil {
		cfg.Previous = locations.Previous
	} else {
		cfg.Previous = nil
	}

	return locations.Results, nil
}

type LocationAreaRes struct {
	Count    int        `json:"count"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []Location `json:"results"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
