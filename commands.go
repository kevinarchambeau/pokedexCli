package main

import (
	"errors"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
	}
}

func commandMap(cfg *config) error {
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

func commandMapb(cfg *config) error {
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

func commandExit(cfg *config) error {
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Print("Available commands: \n\n")
	for _, command := range availableCommands() {
		fmt.Println(command.name + ": " + command.description)
	}

	return nil
}
