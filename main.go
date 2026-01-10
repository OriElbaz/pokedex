package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var commandRegistry = map[string]commandCli {
	"exit" : {
		name: "exit",
		description: "exists program",
		callback: commandExit,
	},
	"help" : {},
}

// init function to avoid circular problem with commandRegistry and commandHelp (looping registry)
func init() {
	commandRegistry["help"] = commandCli{
		name: "help",
		description: "shows users commands",
		callback: commandHelp,
		}
}


func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for i := 0; ; i++ {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		input := scanner.Text()
		
		cleanedInput := cleanInput(input)
		command := strings.ToLower(cleanedInput[0])

		c, err := commandRegistry[command]
		if err == true {
			c.callback()
		} else {
			fmt.Print("Unkown command\n")
		}

	}

}


type commandCli struct {
	name string
	description string
	callback func() error
}


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

