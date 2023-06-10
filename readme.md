# Boule 🎱

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
