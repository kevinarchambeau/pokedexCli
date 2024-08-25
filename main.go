package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	commandList := availableCommands()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Pokedex CLI\n")
	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		if _, ok := commandList[input]; ok {
			err := commandList[input].callback()
			if err != nil {
				fmt.Printf("error occurred: %v", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("invalid command")
		}
	}

}
