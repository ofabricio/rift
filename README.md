# rift

Bind and Unbind values to a struct with a path notation.

## Example

```go
package main

import "github.com/ofabricio/rift"

func main() {

    var user struct {
        Name      string
        Addresses []struct {
            Street string
            Number int
        }
    }

    user.Name = "Luke"

    bs := rift.Bind(&user,
        rift.Field("Name", "John"),
        rift.Field("Addresses.0.Street", "Main"),
        rift.Field("Addresses.0.Number", 100),
        rift.Field("Addresses.1.Street", "Avenue"),
        rift.Field("Addresses.1.Number", 200),
    )

    fmt.Println(user)
    // {John [{Main 100} {Avenue 200}]}

    for _, v := range bs {
        fmt.Println(v.Path, v.Type, v.Old, v.New)
        // Name               string Luke John
        // Addresses.0.Street string      Main
        // Addresses.0.Number int    0    100
        // Addresses.1.Street string      Avenue
        // Addresses.1.Number int    0    200
    }

    for _, v := range rift.Unbind(&user) {
        fmt.Println(v.Path, v.Type, v.Value)
        // Name               string John
        // Addresses.0.Street string Main
        // Addresses.0.Number int    100
        // Addresses.1.Street string Avenue
        // Addresses.1.Number int    200
    }
}
```

## Documentation

### Bind

Bind binds values to a struct based on the provided paths.

It returns a slice of bound values that contain the field's path, type, new value, and old value.

```go
bs := rift.Bind(&user, rift.Field("Addresses.0.Street", "Main"))
```

### Unbind

Unbind extracts public fields from a struct.

It returns a slice of values that contain the field's path, type, and value.

```go
bs := rift.Unbind(&user)
```
