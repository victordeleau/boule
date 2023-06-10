# Boule 🎱

![tests](https://github.com/victordeleau/boule/actions/workflows/test.yaml/badge.svg?event=push)

Boule is a Go boolean expression language. It uses a Context-Free Grammar (CFG) that supports any number of identifiers
of type `STRING`, `NUMBER`, and `BOOLEAN`, as well as recursive expressions using grouping brackets `()`.

Evaluating the expression `!arrived && (origin == "Mars" || (destination == "Titan"))` using the following struct would return `true`:

```go
spaceTravel := &struct{
    Arrived bool
    Origin string
    Destination string
}{
    Arrived: false,
    Origin: "Mars",
    Destination: "Saturn",
}
```

## Example

```go
func main() {
	
	// First create a `boule` expression by passing the expression string.
	// It returns a closure used to evaluate the expression and an error.
	// The expression syntax will be checked against the authorized grammar.
	
    expressionString := "!arrived && (origin == 'Mars' || (destination == 'Titan'))"
    evaluate, err := boule.NewBouleExpression(expressionString)
    if err != nil {
        panic(err)
    }

    // Then, instantiate a prefix tree data structure. Prefix trees enable fast access lookup to the data they contain.
	// To add data, either pass a single `map[string]interface{}` argument, or a key/value pair as two arguments.
    
    data := prefixtree.New()
    
    if err := data.Add(map[string]interface{}{
        "arrived": false,
        "origin":  "Mars",
    }); err != nil {
        panic(err)
    }
    
    if err := data.Add("destination", "Titan"); err != nil {
        panic(err)
    }

    // You can now evaluate the expression against the prefix tree data structure.
    // Call the closure by passing it the data. An error will be returned if type checking failed.
    // Evaluating the expression will either return 'true' or 'false'.
    
    result, err := evaluate(data)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("The expression %q evaluates to '%t'", expressionString, result)
}
```

## Grammar

```
expression         -> binary | suffixExpression
suffixExpression   -> grouping | literal | unary
literal            -> NUMBER | STRING | IDENT
unary              -> NOT suffixExpression
binary             -> expression operator suffixExpression
grouping           -> OPEN expression CLOSE
operator           -> EQUAL | NOT_EQUAL | LESS | LESS_EQUAL | GREATER | GREATER_EQUAL | AND | OR
```
