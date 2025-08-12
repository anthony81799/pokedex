package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
)

func commandHelp(cfg *config, params ...string) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}

func commandCatch(cfg *config, params ...string) error {
	if len(params) == 0 {
		return errors.New("you must provide a pokemon")
	}

	pokemon, err := cfg.pokeapiClient.GetPokemon(params[0])
	if err != nil {
		return err
	}

	res := rand.Intn(pokemon.BaseExperience)

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	if res > 40 {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemon.Name)
	fmt.Println("You may now inspect it with the inspect command.")

	cfg.caughtPokemon[pokemon.Name] = pokemon
	return nil
}

func commandExplore(cfg *config, params ...string) error {
	if len(params) == 0 {
		return errors.New("you must provide a location")
	}

	fmt.Printf("Exploring %s...\n", params[0])
	locationResp, err := cfg.pokeapiClient.GetLocation(params[0])
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, pokemon := range locationResp.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
	return nil
}

func commandInspect(cfg *config, params ...string) error {
	if len(params) == 0 {
		return errors.New("you must provide a pokemon")
	}

	pokemon, ok := cfg.caughtPokemon[params[0]]
	if !ok {
		return errors.New("you have not caught that pokemon")
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, pokemonType := range pokemon.Types {
		fmt.Printf("  - %s\n", pokemonType.Type.Name)
	}

	return nil
}

func commandMapf(cfg *config, params ...string) error {
	locationsResp, err := cfg.pokeapiClient.ListLocations(cfg.nextLocationsURL)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = locationsResp.Next
	cfg.prevLocationsURL = locationsResp.Previous

	for _, loc := range locationsResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(cfg *config, params ...string) error {
	if cfg.prevLocationsURL == nil {
		return errors.New("you're on the first page")
	}

	locationResp, err := cfg.pokeapiClient.ListLocations(cfg.prevLocationsURL)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = locationResp.Next
	cfg.prevLocationsURL = locationResp.Previous

	for _, loc := range locationResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandPokedex(cfg *config, params ...string) error {
	fmt.Println("Your Pokedex:")
	for _, pokemon := range cfg.caughtPokemon {
		fmt.Printf(" - %s\n", pokemon.Name)
	}

	return nil
}

func commandExit(cfg *config, params ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
