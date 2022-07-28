# Deecpy

[![Go Reference](https://img.shields.io/badge/go-reference-%23007d9c?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/unsafe-risk/deecpy)

**Deecpy**, The DeepCopy Library

# Example

```go
package main

import (
    "fmt"

    "github.com/unsafe-risk/deecpy"
)

// func deecpy.Copy[T any](dst, src *T) error
// func deecpy.Duplicate[T any](src T) (T, error)

type Person struct {
    Name string
    Age  int
    id   ID
}

type ID struct {
    UUID string
    email string
}

var john = Person{Name: "John", Age: 30, id: ID{UUID: "123", email: "john@example.com"}}
var jane = Person{Name: "Jane", Age: 25, id: ID{UUID: "456", email: "jane@example.com"}}

func main() {
    var john_copy Person
    err := deecpy.Copy(&john_copy, &john)
    if err != nil {
        panic(err)
    }

    fmt.Println("john:", john)
    fmt.Println("john_copy:", john_copy)

    jane_copy, err := deecpy.Duplicate(jane)
    if err != nil {
        panic(err)
    }

    fmt.Println("jane:", jane)
    fmt.Println("jane_copy:", jane_copy)
}
```
[***Go Playground***](https://go.dev/play/p/ef0QEoCKuTV)
