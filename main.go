package boule

import (
	"fmt"
	"github.com/victordeleau/boule/internal/prefixtree"
)

func main() {

	expression := "!arrived && (origin == \"Mars\" || (destination == \"Titan\"))"

	data := prefixtree.New()
	err := data.Add("arrived", false)
	if err != nil {
		panic(err)
	}
	err = data.Add("origin", "Mars")
	if err != nil {
		panic(err)
	}

	evaluate, err := NewAST(expression, data)
	if err != nil {
		panic(err)
	}

	result, err := evaluate()
	if err != nil {
		panic(err)
	}

	fmt.Printf("The expression '%s' evaluates to '%t'", expression, result)
}
