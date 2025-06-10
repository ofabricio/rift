# rift

Assign values to a struct using path notation to partially update it and track the changes.

## Install

```sh
go get github.com/ofabricio/rift
```

## Examples

This example shows how to assign values to a struct and track the changes.

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

    chg := rift.SetMany(&user,
        rift.Path("Name", "John"),
        rift.Path("Addresses.0.Street", "Main"),
        rift.Path("Addresses.0.Number", 100),
        rift.Path("Addresses.1.Street", "Avenue"),
        rift.Path("Addresses.1.Number", 200),
    )

    fmt.Println(user)
    // {John [{Main 100} {Avenue 200}]}

    for _, v := range chg {
        fmt.Println(v.Path, v.Type, v.Old, v.New)
        // Name               string Luke John
        // Addresses.0.Street string      Main
        // Addresses.0.Number int    0    100
        // Addresses.1.Street string      Avenue
        // Addresses.1.Number int    0    200
    }
}
```

Note that `Set*` functions return the changes.

### Tree

This example shows how to get a tree representation of the provided struct.

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

    // Other way to assign values.
    rift.SetPath(&user, "Name", "John")
    rift.SetPath(&user, "Addresses.0.Street", "Main")
    rift.SetPath(&user, "Addresses.0.Number", 100)
    rift.SetPath(&user, "Addresses.1.Street", "Avenue")
    rift.SetPath(&user, "Addresses.1.Number", 200)

    tree := rift.Get(user)

    data, _ := json.MarshalIndent(tree, "", "    ")

    fmt.Println(string(data))

    // Output:
    // {
    //     "Name": "",
    //     "Path": "",
    //     "Type": "struct",
    //     "Data": null,
    //     "Next": [
    //         {
    //             "Name": "Name",
    //             "Path": "Name",
    //             "Type": "string",
    //             "Data": "John",
    //             "Next": null
    //         },
    //         {
    //             "Name": "Addresses",
    //             "Path": "Addresses",
    //             "Type": "slice",
    //             "Data": null,
    //             "Next": [
    //                 {
    //                     "Name": "0",
    //                     "Path": "Addresses.0",
    //                     "Type": "struct",
    //                     "Data": null,
    //                     "Next": [
    //                         {
    //                             "Name": "Street",
    //                             "Path": "Addresses.0.Street",
    //                             "Type": "string",
    //                             "Data": "Main",
    //                             "Next": null
    //                         },
    //                         {
    //                             "Name": "Number",
    //                             "Path": "Addresses.0.Number",
    //                             "Type": "int",
    //                             "Data": 100,
    //                             "Next": null
    //                         }
    //                     ]
    //                 },
    //                 {
    //                     "Name": "1",
    //                     "Path": "Addresses.1",
    //                     "Type": "struct",
    //                     "Data": null,
    //                     "Next": [
    //                         {
    //                             "Name": "Street",
    //                             "Path": "Addresses.1.Street",
    //                             "Type": "string",
    //                             "Data": "Avenue",
    //                             "Next": null
    //                         },
    //                         {
    //                             "Name": "Number",
    //                             "Path": "Addresses.1.Number",
    //                             "Type": "int",
    //                             "Data": 200,
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

Now it is possible to update one of those nodes in the tree and
pass it to `rift.Set(dst any, n Node)` function to update just that part of the struct.

Example:

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

    node := rift.Node{
        Next: []rift.Node{
            {
                Path: "Name",
                Data: "John",
            },
            {
                Next: []rift.Node{
                    {
                        Path: "Addresses.0.Street",
                        Data: "Main",
                    },
                },
            },
        },
    }

    rift.Set(&user, node)

    fmt.Println(user)

    // Output:
    // {John [{Main 0}]}
}
```

Note that even though `Node` has more informations, only `Path` and `Data` are required to `Set`.
Also only nodes with `Next == nil` are applied.
