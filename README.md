# Vow

An attempt at a valibot-like validation library for Go.

## Installation

```bash
go get github.com/smc13/vow
```

## Example Usage

```go
package main

import (
  "fmt"
  "github.com/smc13/vow"
)

type User struct {
  Name  string `vow:"name"` // optional struct tag to map fields
  Age   int `vow:"age"`
  Email string `vow:"email"`
}

func main() {
  schema := vow.Struct[User](vow.Fields{
    "name": Pipe(String(), MinLen(1), MaxLen(100)),
    "age":  Pipe(Number(), Min(0), Max(150)),
    "email": Pipe(String(), Email()),
  })

  data := map[string]interface{}{
    "name":  "Alice",
    "age":   30,
    "email": "test@example.com"
  }

  var user User
  err := vow.Parse(context.Background(), schema, data, &user) // or vow.ParseAs for generic types
  
  // fully typed and validated struct
  fmt.Printf("%+v\n", user)
}
```

## Available Schemas
  - [x] Struct
  - [ ] Map
  - [ ] Slice
  - [x] String
    - [ ] Email
  - [x] Int, Float
    - [ ] Max
    - [ ] Min
  - [x] Bool
