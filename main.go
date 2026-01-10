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

		c, err := commandRegistry[command]
		if err == true {
			c.callback()
		} else {
			fmt.Print("Unkown command\n")
		}

	}

}


