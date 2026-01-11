package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	
	for i := 0; ; i++ {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		input := scanner.Text()
		
		cleanedInput := cleanInput(input)
		command := strings.ToLower(cleanedInput[0])
		

		commandInput := ""
		
		if len(cleanedInput) > 1 {
			commandInput = cleanedInput[1]	
		}


		c, err := commandRegistry[command]
		if err == true {
			c.callback(commandInput)
		} else {
			fmt.Print("Unkown command\n")
		}

	}

}


