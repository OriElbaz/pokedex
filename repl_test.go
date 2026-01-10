package main

import (
	"reflect"
	"testing"
	"fmt"
)


func TestCleanInput(t *testing.T) {
	tests := []struct {
		name 		string
		input 		string
		expected 	[]string
	}{
	{name: "1", input: "   Hello  world  ", expected: []string{"Hello", "world"}},
	{name: "2", input: "helloWorld", expected: []string{"helloWorld"}},
	{name: "3", input: "", expected: []string{}},
	}
	
	for _, test := range tests {
		got := cleanInput(test.input)
		if !reflect.DeepEqual(got, test.expected) {
			t.Errorf("Test: %s, Expected: %s, Got: %s", test.name, test.expected, got)
		} else {
			fmt.Printf("Passed test: %s, Expected: %s, Got: %s\n", test.name, test.expected, got)
		}

	}
}