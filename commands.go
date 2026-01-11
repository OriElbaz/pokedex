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
	"explore" : {
		name: "explore",
		description: "shows pokemon in area",
		callback: explore,
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
	callback func(string) error
	config *map[string]string
}

// COMMAND CALLBACK FUNCTIONS

func commandExit(none string) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(none string) error {
	fmt.Print("Welcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n")

	for key, command := range commandRegistry {
		fmt.Printf("%s: %s\n", key, command.description)
	}

	return nil
}

func commandMap(none string) error {
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
		
		if res.StatusCode != http.StatusOK {
			fmt.Printf("Unexpected status code: %d", res.StatusCode)
			return nil
		}

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

func commandMapb(none string) error {
	if AfterId > 0 {
		AfterId -= 40
	} else if AfterId == 0 {
		AfterId = -20
	}

	err := commandMap(none)

	return err
}

type pokemonInLocation struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func explore(location string) error {
	
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", location)
	var pokemons pokemonInLocation

	data, ok := cache.Get(url)
	if !ok {
		// http get request //
		res, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error with GET: %v\n", err)
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			fmt.Printf("Unexpected status code: %d", res.StatusCode)
			return nil
		}

		// cache data //
		data, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error with READALL: %v\n", err)
			return err
		}
		cache.Add(url, data)
	}

	// convert byte into struct //
	err := json.Unmarshal(data, &pokemons)
	if err != nil {
		fmt.Printf("Error with UNMARSHAL: %v\n", err)
		return err
	}

	for _, pokemon := range pokemons.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}

	return nil
}