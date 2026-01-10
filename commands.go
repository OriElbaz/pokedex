package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)


var AfterId = -20
var ptr *int = &AfterId

var urls = map[string]string{
	"Next": "https://pokeapi.co/api/v2/location-area/?limit=20&offset=",
	"Previous": fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?limit=20&offset=%d", *ptr),
}

var commandRegistry = map[string]commandCli {
	"exit" : {
		name: "exit",
		description: "exists program",
		callback: commandExit,
		config: &urls,
	},
	"help" : {},
	"map" : {
		name: "map",
		description: "shows 20 locations",
		callback: commandMap,
		config: &urls,
	},
	"mapb" : {
		name: "map",
		description: "shows previous 20 locations",
		callback: commandMapb,
		config: &urls,
	},
}

// init function to avoid circular problem with commandRegistry and commandHelp (looping registry)
func init() {
	commandRegistry["help"] = commandCli{
		name: "help",
		description: "shows users commands",
		callback: commandHelp,
		config: &urls,
		}
}

type commandMapStruct struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type commandCli struct {
	name string
	description string
	callback func() error
	config *map[string]string
}

// COMMAND CALLBACK FUNCTIONS

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Print("Welcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n")

	for key, command := range commandRegistry {
		fmt.Printf("%s: %s\n", key, command.description)
	}

	return nil
}

func commandMap() error {
	AfterId += 20

	url_ := "https://pokeapi.co/api/v2/location-area/?limit=20&offset="
	urls["Next"] = url_ + fmt.Sprintf("%d", AfterId)

	/*
	PRINTING TO ENSURE CORRECTNESS:
	fmt.Printf("AFTER ID: %d\n", AfterId)
	fmt.Printf("URL: %s\n\n", urls["Next"])
	*/
	url := fmt.Sprint(urls["Next"])

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return err
	}
	defer res.Body.Close()

	// data, err := io.ReadAll(res.Body)
	// fmt.Print(string(data))

	var locations commandMapStruct
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&locations); err != nil {
		fmt.Printf("Error: %v", err)
		return err
	}

	for _, location := range locations.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}

func commandMapb() error {
	if AfterId > 0 {
		AfterId -= 40
	} else if AfterId == 0 {
		AfterId = -20
	}

	err := commandMap()
	
	if AfterId == 40 {
		AfterId = 20
	}

	return err
}