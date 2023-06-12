package main

import (
	"fmt"
	"github.com/victordeleau/boule"
)

func main() {

	// First create a `boule` expression by passing the expression string.
	// The expression syntax will be checked against the authorized grammar.
	expressionString := "!arrived && (origin == 'Mars' || (destination == 'Titan')) && !cancelled"
	evaluate, _ := boule.NewExpression(expressionString)

	// Then, instantiate a boule.Data struct. It uses a prefix tree for fast access lookup.
	data := boule.NewData()

	// structs are supported
	_ = data.Add(struct {
		Cancelled bool `json:"cancelled"`
	}{
		Cancelled: false,
	})

	// maps are supported
	_ = data.Add(map[string]interface{}{
		"arrived": false,
		"origin":  "Mars",
	})

	// key/value pairs are supported
	_ = data.Add("destination", "Titan")

	// You can now evaluate the expression against the data, which returns 'true' or 'false'.
	// An error will be returned if type checking failed, or if an identifier was not found.
	result, _ := evaluate(data)

	fmt.Printf("The expression %q evaluates to '%t'", expressionString, result)
}
