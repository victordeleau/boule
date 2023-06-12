# Boule 🎱

![tests](https://github.com/victordeleau/boule/actions/workflows/test.yaml/badge.svg?event=push)

Boule is a Go boolean expression language. It uses a Context-Free Grammar (CFG) that supports any number of identifiers
of type `STRING`, `NUMBER`, and `BOOLEAN`, as well as recursive expressions using grouping brackets `()`.

Expressions are evaluated against a prefix-tree data structure containing the identifiers in the expression.
Data can be loaded into the prefix-tree as key/value pairs, as a `map[string]interface{}`, or as a struct.

For structs passed as data, any number of embedded structs are supported, but not maps nor slices (yet).
The identifier name for structs is the json name of the field, which is required for the field to be considered.

## Example

```go
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
