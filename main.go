package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type config struct {
	apiClient     client
	nextLocations *string
	prevLocations *string
}

func main() {
	appConfig := config{
		apiClient: newClient(time.Minute * 5),
	}

	commandList := availableCommands()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Pokedex CLI\n\n")
	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		args := []string{}
		words := cleanInput(input)
		command := words[0]
		if len(input) > 1 {
			args = words[1:]
		}

		if _, ok := commandList[command]; ok {
			err := commandList[command].callback(&appConfig, args...)
			if err != nil {
				fmt.Printf("error occurred: %v", err)
			}
		} else {
			fmt.Println("invalid command")
		}
	}

}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
