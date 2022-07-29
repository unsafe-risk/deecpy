[![Go Reference](https://img.shields.io/badge/go-reference-%23007d9c?style=for-the-badge&logo=go)](https://pkg.go.dev/v8.run/go/deecpy)
[![GitHub Workflow Status (push)](https://img.shields.io/github/workflow/status/unsafe-risk/deecpy/go_test/main?event=push&style=for-the-badge)](https://github.com/unsafe-risk/deecpy/actions/workflows/go_test.yml)
[![Codecov main](https://img.shields.io/codecov/c/gh/unsafe-risk/deecpy/main?style=for-the-badge)](https://app.codecov.io/gh/unsafe-risk/deecpy)
[![GitHub license](https://img.shields.io/github/license/unsafe-risk/deecpy?style=for-the-badge)](https://pkg.go.dev/v8.run/go/deecpy?tab=licenses)
# Deecpy

**Deecpy**, The DeepCopy Library

# Example

```go
package main

import (
    "fmt"

    "v8.run/go/deecpy"
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
[***Go Playground***](https://go.dev/play/p/4Kc5gPCaebk)
