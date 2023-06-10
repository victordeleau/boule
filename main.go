package boule

import (
	"fmt"
	"github.com/victordeleau/boule/internal/prefixtree"
)

func main() {

	expression := "!arrived && (origin == \"Mars\" || (destination == \"Titan\"))"

	data := prefixtree.New()
	if err := data.Add("arrived", false); err != nil {
		panic(err)
	}

	if err := data.Add("origin", "Mars"); err != nil {
		panic(err)
	}

	evaluate, err := NewBouleExpression(expression, data)
	if err != nil {
		panic(err)
	}

	result, err := evaluate()
	if err != nil {
		panic(err)
	}

	fmt.Printf("The expression '%s' evaluates to '%t'", expression, result)
}
