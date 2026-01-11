package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/OriElbaz/pokedex/internal/pokecache"
)


var AfterId = -20
var cache = *pokecache.NewCache(30 * time.Second)

var urls = map[string]string{
	"Next": "https://pokeapi.co/api/v2/location-area/?limit=20&offset=",
	"Previous": "https://pokeapi.co/api/v2/location-area/?limit=20&offset=",
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
	if AfterId == 0 {
		urls["Previous"] = url_ + fmt.Sprintf("%d", AfterId)
	} else {
		urls["Previous"] = url_ + fmt.Sprintf("%d", AfterId-20)
	}
	
	urlToUse := fmt.Sprint(urls["Next"])
	data := []byte{}

	entry, ok := cache.Get(urlToUse); 
	// fmt.Printf("***** USED CACHE: %v\n", ok)

	if !ok {

		res, err := http.Get(urlToUse)
		if err != nil {
			fmt.Printf("Error with GET: %v\n", err)
			return err
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error with READALL: %v\n", err)
			return err
		}

		cache.Add(urlToUse, data)
		// fmt.Printf("***** CACHED: %v\n", urlToUse)

	} else {
		data = entry
	}

	decodeBody := bytes.NewReader(data)

	var locations commandMapStruct
	decoder := json.NewDecoder(decodeBody)
	if err := decoder.Decode(&locations); err != nil {
		fmt.Printf("Error with DECODING: %v\n", err)
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

	return err
}