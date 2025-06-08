# rift

Bind and Unbind values to a struct with a path notation.

## Bind and Unbind Example

This example shows how `Bind` function binds values to a struct based on the provided paths; and shows
how `Unbind` function extracts public fields from a struct.

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

Note that `Bind` function returns a slice of bound values, which contains the field path, type, old value and new value.

Note that `Unbind` function returns a slice of unbound values, which contains the field path, type and value.

## Describe Example

This example shows how `Describe` function returns a tree representation of the provided struct.

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

    rift.Bind(&user,
        rift.Field("Name", "John"),
        rift.Field("Addresses.0.Street", "Main"),
        rift.Field("Addresses.0.Number", 100),
        rift.Field("Addresses.1.Street", "Avenue"),
        rift.Field("Addresses.1.Number", 200),
    )

    tree := rift.Describe(user)

    data, _ := json.MarshalIndent(tree, "", "    ")

    fmt.Println(string(data))

    // Output:
    // {
    //     "Name": "",
    //     "Path": "",
    //     "Type": "struct",
    //     "Value": null,
    //     "Next": [
    //         {
    //             "Name": "Name",
    //             "Path": "Name",
    //             "Type": "string",
    //             "Value": "John",
    //             "Next": null
    //         },
    //         {
    //             "Name": "Addresses",
    //             "Path": "Addresses",
    //             "Type": "slice",
    //             "Value": null,
    //             "Next": [
    //                 {
    //                     "Name": "0",
    //                     "Path": "Addresses.0",
    //                     "Type": "struct",
    //                     "Value": null,
    //                     "Next": [
    //                         {
    //                             "Name": "Street",
    //                             "Path": "Addresses.0.Street",
    //                             "Type": "string",
    //                             "Value": "Main",
    //                             "Next": null
    //                         },
    //                         {
    //                             "Name": "Number",
    //                             "Path": "Addresses.0.Number",
    //                             "Type": "int",
    //                             "Value": 100,
    //                             "Next": null
    //                         }
    //                     ]
    //                 },
    //                 {
    //                     "Name": "1",
    //                     "Path": "Addresses.1",
    //                     "Type": "struct",
    //                     "Value": null,
    //                     "Next": [
    //                         {
    //                             "Name": "Street",
    //                             "Path": "Addresses.1.Street",
    //                             "Type": "string",
    //                             "Value": "Avenue",
    //                             "Next": null
    //                         },
    //                         {
    //                             "Name": "Number",
    //                             "Path": "Addresses.1.Number",
    //                             "Type": "int",
    //                             "Value": 200,
    //                             "Next": null
    //                         }
    //                     ]
    //                 }
    //             ]
    //         }
    //     ]
    // }
}
```
