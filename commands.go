package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/OriElbaz/pokedex/internal/pokecache"
)

/*** VARIABLES ***/
var AfterId = -20

var cache = *pokecache.NewCache(30 * time.Second)

var pokedex = make(map[string]Pokemon)

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
		callback: commandExplore,
		config: &urls,
	},
	"catch" : {
		name: "catch",
		description: "attempt to catch a pokemon",
		callback: commandCatch,
		config: &urls,
	},
	"inspect" : {
		name: "inspect",
		description: "shows a pokemon's, you've caught, stats",
		callback: commandInspect,
		config: &urls,
	},
	"pokedex" : {
		name: "pokedex",
		description: "shows all pokemon in your pokedex",
		callback: commandPokedex,
		config: &urls,
	},
}



/*** COMMAND STRUCTS ***/
// map
type commandMapStruct struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
// all
type commandCli struct {
	name string
	description string
	callback func(string) error
	config *map[string]string
}
// explore
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

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}


/*** COMMAND CALLBACK FUNCTIONS ***/

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

	urlToUse := updateUrlsInMap()
	var locations commandMapStruct
	data, ok := cache.Get(urlToUse); 

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
	}

	err := json.Unmarshal(data, &locations)
	if err != nil {
		fmt.Printf("Error with UNMARSHAL: %v\n", err)
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

func commandExplore(location string) error {
	
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

func commandCatch(pokemonInput string) error {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonInput
	var pokemon Pokemon

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

	// read into bytes //
	pokemonData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error with READALL: %v\n", err)
		return err
	}

	// convert byte into struct //
	err = json.Unmarshal(pokemonData, &pokemon)
	if err != nil {
		fmt.Printf("Error with UNMARSHAL: %v\n", err)
		return err
	}
	
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	if !catchChance(pokemon.BaseExperience) {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemon.Name)
	pokedex[pokemon.Name] = pokemon
	
	
	return nil

}

func commandInspect(pokemonInput string) error {
	pokemon, err := pokedex[pokemonInput]
	if err == false {
		fmt.Printf("you have not caught that pokemon\n")
		return nil
	}

	printStats(pokemon)

	return nil
}

func commandPokedex(none string) error {
	fmt.Println("Your Pokedex:")
	for _, pokemon := range pokedex {
		fmt.Printf("- %s\n", pokemon.Name)
	}
	return nil
}

/*** HELPER FUNCTIONS ***/
func init() {
	/* 
	Function is used to add the 'help' commandCli into commandRegistry after commandRegistry is instantiated.
	This stops circular reliance on each other.
	*/

	commandRegistry["help"] = commandCli{
		name: "help",
		description: "shows users commands",
		callback: commandHelp,
		config: &urls,
		}
}

func updateUrlsInMap() string {
	url_ := "https://pokeapi.co/api/v2/location-area/?limit=20&offset="
	urls["Next"] = url_ + fmt.Sprintf("%d", AfterId)
	if AfterId == 0 {
		urls["Previous"] = url_ + fmt.Sprintf("%d", AfterId)
	} else {
		urls["Previous"] = url_ + fmt.Sprintf("%d", AfterId-20)
	}
	return fmt.Sprint(urls["Next"])
}

func catchChance(baseExperience int) bool {
	chance := 100.0 - (float64(baseExperience) - 30.0) * (90.0 / 578.0)
	if rand.Float64() < chance {
		return true
	}
	return false

}

func printStats(pokemon Pokemon) {
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		name := stat.Stat.Name
		num := stat.BaseStat
		fmt.Printf(" - %s: %d\n", name, num)
	}
	fmt.Printf("Types:\n")
	for _, t := range pokemon.Types {
		typeName := t.Type.Name
		fmt.Printf(" - %s\n", typeName)
	}
}