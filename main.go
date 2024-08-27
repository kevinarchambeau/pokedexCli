package main

import (
	"bufio"
	"fmt"
	"os"
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
