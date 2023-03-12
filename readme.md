# Boule 🎱

Boule is a Go boolean expression language. It uses a Context-Free Grammar (CFG) that supports any number of identifiers
of type `string`, `integer`, `float` and `boolean`

Evaluating the expression `!arrived && (origin == "Mars")` using the struct

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

would return `true`.