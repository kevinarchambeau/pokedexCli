package main

import (
	"errors"
	"fmt"
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
	}
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
		fmt.Printf(" - %s\n", enc.Pokemon.Name)
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
