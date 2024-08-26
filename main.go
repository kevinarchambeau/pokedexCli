package main

import (
	"bufio"
	"fmt"
	"os"
)

type config struct {
	apiClient     client
	nextLocations *string
	prevLocations *string
}

func main() {
	appConfig := config{
		apiClient: newClient(),
	}

	commandList := availableCommands()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Pokedex CLI\n")
	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		if _, ok := commandList[input]; ok {
			err := commandList[input].callback(&appConfig)
			if err != nil {
				fmt.Printf("error occurred: %v", err)
			}
		} else {
			fmt.Println("invalid command")
		}
	}

}
