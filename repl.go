package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		words := cleanInput(userInput)
		if len(words) == 0 {
			fmt.Print("Please enter a command\n")
			continue
		}
		firstWord := words[0]
		args := words[1:]
		command, ok := cfg.commands[firstWord]
		if !ok {
			fmt.Print("Unknown command\n")
		} else {
			err := command.callback(cfg, args...)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}
}

func cleanInput(text string) []string {
	lText := strings.ToLower(text)
	var slice []string
	slice = strings.Fields(lText)
	return slice
}
