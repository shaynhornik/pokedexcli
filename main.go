package main

import (
	"fmt"
	"os"
	"io"
	"log"
	"net/http"
	"encoding/json"
)

type cliCommand struct {
	name string
	description string
	callback func(*config) error
}

type config struct {
	commands map[string]cliCommand
	previousUrl *string
	nextUrl *string
}

type locationsPage struct {
	Count int `json:"count"`
	Next *string `json:"next"`
	Previous *string `json:"previous"`
	Results []Location `json:"results"`
}

type Location struct {
	Name string `json:"name"`
	Url string `json:"url"`
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
	res, err := http.Get(*cfg.nextUrl)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	locations := locationsPage{}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(locations.Results); i++ {
		fmt.Printf("%s\n", locations.Results[i].Name)
	}
	cfg.previousUrl = locations.Previous
	if locations.Next != nil{
		cfg.nextUrl = locations.Next
	} else {
		fmt.Print("No more locations\n")
	}
	return nil
}

func commandMapb(cfg *config) error {
	if cfg.previousUrl == nil {
		fmt.Print("You're on the first page, no previous locations\n")
		return nil
	}
	res, err := http.Get(*cfg.previousUrl)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	locations := locationsPage{}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(locations.Results); i++ {
		fmt.Printf("%s\n", locations.Results[i].Name)
	}
	cfg.nextUrl = locations.Next
	cfg.previousUrl = locations.Previous
	return nil
}

var baseUrl string = "https://pokeapi.co/api/v2/location-area"

func main() {

	commands := map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"map": {
			name: "map",
			description: "Displays the names of 20 location areas",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "Displays the names of the previous 20 location areas",
			callback: commandMapb,
		},
	}

	cfg := config{
		commands: commands,
		previousUrl: nil,
		nextUrl: &baseUrl,
	}

	startRepl(&cfg)
}
