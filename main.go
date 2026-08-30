package main

import (
	"fmt"
	"os"
	"time"

	"pokedexcli/internal/pokeapi"
	"pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	commands      map[string]cliCommand
	pokeapiClient pokeapi.Client
	nextURL       *string
	previousURL   *string
}

func commandExit(cfg *config) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex\n")
	return nil
}

func commandMap(cfg *config) error {
	locations, err := cfg.pokeapiClient.ListLocations(cfg.nextURL)
	if err != nil {
		return err
	}

	for _, loc := range locations.Results {
		fmt.Printf("%s\n", loc.Name)
	}

	cfg.previousURL = locations.Previous
	cfg.nextURL = locations.Next

	return nil
}

func commandMapb(cfg *config) error {
	if cfg.previousURL == nil {
		fmt.Print("You're on the first page, no previous locations\n")
		return nil
	}

	locations, err := cfg.pokeapiClient.ListLocations(cfg.previousURL)
	if err != nil {
		return err
	}

	for _, loc := range locations.Results {
		fmt.Printf("%s\n", loc.Name)
	}

	cfg.nextURL = locations.Next
	cfg.previousURL = locations.Previous

	return nil
}

func main() {
	interval := 5 * time.Second
	c := pokecache.NewCache(interval)
	pokeClient := pokeapi.NewClient(interval, c)

	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the previous 20 location areas",
			callback:    commandMapb,
		},
	}

	cfg := config{
		commands:      commands,
		pokeapiClient: pokeClient,
	}

	startRepl(&cfg)
}
