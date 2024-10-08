package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func availableCommands() map[string]cliCommand {
	return map[string]cliCommand{
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
			description: "Displays a list of 20 locations, repeating returns new sets",
			callback:    commandMap,
		},
		"mapb": {
			name:        "map",
			description: "Displays a list of previous 20 locations, repeating works back to the beginning",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore <location_name>",
			description: "Explore a location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon>",
			description: "attempts to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon>",
			description: "gives pokemon stats (if caught)",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "prints all pokemon in your pokedex",
			callback:    commandPokedex,
		},
	}
}

func commandPokedex(cfg *config, args ...string) error {

	if len(cfg.caught) == 0 {
		fmt.Println("You haven't caught any pokemon yet")
		return nil
	}
	fmt.Println("Pokedex contents:")
	for key := range cfg.caught {
		fmt.Printf(" %s\n", key)
	}

	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("no pokemon provided")
	}

	name := args[0]

	pokemon, ok := cfg.caught[name]
	if !ok {
		fmt.Printf("You don't have a %s\n", name)
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Println("Stats: ")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  %s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Type/s:")
	for _, pokeType := range pokemon.Types {
		fmt.Printf("  %s\n", pokeType.Type.Name)
	}

	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("no pokemon provided")
	}

	name := args[0]
	pokemon, err := cfg.apiClient.GetPokemonInfo(name)
	if err != nil {
		return err
	}

	if _, ok := cfg.caught[pokemon.Name]; ok {
		fmt.Printf("You've already caught a %s\n", pokemon.Name)
		return nil
	}

	chance := rand.Intn(pokemon.BaseExperience)
	cfg.caught[pokemon.Name] = pokemon
	if chance >= pokemon.BaseExperience/2 {
		fmt.Printf("You caught a %s\n", pokemon.Name)
		cfg.caught[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s ran away\n", pokemon.Name)
	}

	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("no location provided")
	}

	name := args[0]
	location, err := cfg.apiClient.GetLocation(name)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s\n", location.Name)
	fmt.Println("Found pokemon:")
	for _, enc := range location.PokemonEncounters {
		fmt.Printf(" %s\n", enc.Pokemon.Name)
	}

	return nil
}

func commandMap(cfg *config, args ...string) error {
	response, err := cfg.apiClient.ListLocations(cfg.nextLocations)
	if err != nil {
		return err
	}

	cfg.nextLocations = response.Next
	cfg.prevLocations = response.Previous

	for _, location := range response.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.prevLocations == nil {
		return errors.New("no previous locations to go to\n")
	}

	response, err := cfg.apiClient.ListLocations(cfg.prevLocations)
	if err != nil {
		return err
	}

	cfg.nextLocations = response.Next
	cfg.prevLocations = response.Previous

	for _, location := range response.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExit(cfg *config, args ...string) error {
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Print("Available commands: \n\n")
	for _, command := range availableCommands() {
		fmt.Println(command.name + ": " + command.description)
	}

	return nil
}
