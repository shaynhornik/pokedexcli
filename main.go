package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"pokedexcli/internal/pokeapi"
	"pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

type config struct {
	commands      map[string]cliCommand
	pokeapiClient pokeapi.Client
	nextURL       *string
	previousURL   *string
}

func commandExit(cfg *config, args ...string) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex\n")
	return nil
}

func commandMap(cfg *config, args ...string) error {
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

func commandMapb(cfg *config, args ...string) error {
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

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("you must provide a location area name")
	}
	areaName := args[0]

	fmt.Printf("Exploring %s...\n", areaName)

	locationArea, err := cfg.pokeapiClient.GetLocationArea(areaName)
	if err != nil {
		return err
	}

	fmt.Print("Found Pokemon:\n")
	for _, encounter := range locationArea.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func main() {
	cacheTTL := 5 * time.Minute
	httpTimeout := 5 * time.Second
	c := pokecache.NewCache(cacheTTL)
	pokeClient := pokeapi.NewClient(httpTimeout, c)

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
		"explore": {
			name:        "explore",
			description: "Takes an input of area name and returns the name of the pokempn found there",
			callback:    commandExplore,
		},
	}

	cfg := config{
		commands:      commands,
		pokeapiClient: pokeClient,
	}

	startRepl(&cfg)
}
