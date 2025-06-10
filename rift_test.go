package rift_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/ofabricio/rift"
)

func Example() {

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

	fmt.Println("User:")
	fmt.Println(user)

	fmt.Println("Changes:")
	for _, v := range chg {
		fmt.Println(v.Path, v.Type, v.Old, v.New)
	}

	// Output:
	// User:
	// {John [{Main 100} {Avenue 200}]}
	// Changes:
	// Name string Luke John
	// Addresses.0.Street string  Main
	// Addresses.0.Number int 0 100
	// Addresses.1.Street string  Avenue
	// Addresses.1.Number int 0 200
}

func ExampleGet() {

	var user struct {
		Name      string
		Addresses []struct {
			Street string
			Number int
		}
	}

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

func ExampleSet() {

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

func TestSetMany(t *testing.T) {

	tt := []struct {
		Desc string
		Give any
		When []rift.Node
		Then any
		Bnds []rift.Change
	}{
		{
			Desc: "given a nil source, it should set the root if field is empty",
			Give: nil,
			When: []rift.Node{
				rift.Path("", 3),
			},
			Then: 3,
			Bnds: []rift.Change{{Path: "", Type: "int", New: 3, Old: nil}},
		},
		{
			Desc: "given a nil source, it should set the root as a slice if field is number",
			Give: nil,
			When: []rift.Node{
				rift.Path("0", 3),
			},
			Then: []any{3},
			Bnds: []rift.Change{{Path: "0", Type: "int", New: 3, Old: nil}},
		},
		{
			Desc: "given a nil source, it should set the root as a slice of len=2 if field is number 1",
			Give: nil,
			When: []rift.Node{
				rift.Path("1", 3),
			},
			Then: []any{nil, 3},
			Bnds: []rift.Change{{Path: "1", Type: "int", New: 3, Old: nil}},
		},
		{
			Desc: "given a nil source, it should set the root as a slice of len=2",
			Give: nil,
			When: []rift.Node{
				rift.Path("0", 2),
				rift.Path("1", 3),
			},
			Then: []any{2, 3},
			Bnds: []rift.Change{
				{Path: "0", Type: "int", New: 2, Old: nil},
				{Path: "1", Type: "int", New: 3, Old: nil},
			},
		},
		{
			Desc: "test inverted slice indexes",
			Give: nil,
			When: []rift.Node{
				rift.Path("1", 2),
				rift.Path("0", 3),
			},
			Then: []any{3, 2},
			Bnds: []rift.Change{
				{Path: "1", Type: "int", New: 2, Old: nil},
				{Path: "0", Type: "int", New: 3, Old: nil},
			},
		},
		{
			Desc: "given a nil source, it should set the root as a map with key a",
			Give: nil,
			When: []rift.Node{
				rift.Path("a", 2),
			},
			Then: map[string]any{"a": 2},
			Bnds: []rift.Change{
				{Path: "a", Type: "int", New: 2, Old: nil},
			},
		},
		{
			Desc: "given a nil source, it should set the root as a map with key a and b",
			Give: nil,
			When: []rift.Node{
				rift.Path("a", 2),
				rift.Path("b", 3),
			},
			Then: map[string]any{"a": 2, "b": 3},
			Bnds: []rift.Change{
				{Path: "a", Type: "int", New: 2, Old: nil},
				{Path: "b", Type: "int", New: 3, Old: nil},
			},
		},
		{
			Desc: "given a nil source, it should set the root as a map with key a and b; and key c should also be a map",
			Give: nil,
			When: []rift.Node{
				rift.Path("a", 2),
				rift.Path("b", 3),
				rift.Path("c.a", 3),
			},
			Then: map[string]any{"a": 2, "b": 3, "c": map[string]any{"a": 3}},
			Bnds: []rift.Change{
				{Path: "a", Type: "int", New: 2, Old: nil},
				{Path: "b", Type: "int", New: 3, Old: nil},
				{Path: "c.a", Type: "int", New: 3, Old: nil},
			},
		},
		{
			Desc: "given a nil source, it should set the root as a map with key a and b; and key c should also be a map with key a and b",
			Give: nil,
			When: []rift.Node{
				rift.Path("a", 2),
				rift.Path("b", 3),
				rift.Path("c.a", 4),
				rift.Path("c.b", 5),
			},
			Then: map[string]any{"a": 2, "b": 3, "c": map[string]any{"a": 4, "b": 5}},
			Bnds: []rift.Change{
				{Path: "a", Type: "int", New: 2, Old: nil},
				{Path: "b", Type: "int", New: 3, Old: nil},
				{Path: "c.a", Type: "int", New: 4, Old: nil},
				{Path: "c.b", Type: "int", New: 5, Old: nil},
			},
		},
		{
			Desc: "if a subfield of a generic map key has an index, that field should be a slice",
			Give: nil,
			When: []rift.Node{
				rift.Path("a.b.0", 3),
			},
			Then: map[string]any{"a": map[string]any{"b": []any{3}}},
			Bnds: []rift.Change{
				{Path: "a.b.0", Type: "int", New: 3, Old: nil},
			},
		},
		{
			Desc: "same as above, but with index 1",
			Give: nil,
			When: []rift.Node{
				rift.Path("a.b.1", 3),
			},
			Then: map[string]any{"a": map[string]any{"b": []any{nil, 3}}},
			Bnds: []rift.Change{
				{Path: "a.b.1", Type: "int", New: 3, Old: nil},
			},
		},
		{
			Desc: "same as above, but with index 0 and 1",
			Give: nil,
			When: []rift.Node{
				rift.Path("a.a.0", 2),
				rift.Path("a.a.1", 3),
			},
			Then: map[string]any{"a": map[string]any{"a": []any{2, 3}}},
			Bnds: []rift.Change{
				{Path: "a.a.0", Type: "int", New: 2, Old: nil},
				{Path: "a.a.1", Type: "int", New: 3, Old: nil},
			},
		},
		{
			Desc: "same as above, but with inverted indexes",
			Give: nil,
			When: []rift.Node{
				rift.Path("a.a.1", 2),
				rift.Path("a.a.0", 3),
			},
			Then: map[string]any{"a": map[string]any{"a": []any{3, 2}}},
			Bnds: []rift.Change{
				{Path: "a.a.1", Type: "int", New: 2, Old: nil},
				{Path: "a.a.0", Type: "int", New: 3, Old: nil},
			},
		},
		{
			Desc: "now a slice with a map inside",
			Give: nil,
			When: []rift.Node{
				rift.Path("a.a.0", 2),
				rift.Path("a.a.1.a", 3),
				rift.Path("a.a.1.b", 4),
			},
			Then: map[string]any{"a": map[string]any{"a": []any{2, map[string]any{"a": 3, "b": 4}}}},
			Bnds: []rift.Change{
				{Path: "a.a.0", Type: "int", New: 2, Old: nil},
				{Path: "a.a.1.a", Type: "int", New: 3, Old: nil},
				{Path: "a.a.1.b", Type: "int", New: 4, Old: nil},
			},
		},
		{
			Desc: "test setting many different field types of a struct",
			Give: &TestData{},
			When: []rift.Node{
				rift.Path("Int", 2),
				rift.Path("IntPtr", 3),
				rift.Path("String", "Hi"),
				rift.Path("Slice.0.Int", 1),
				rift.Path("SlicePtr.0.Int", 1),
				rift.Path("Struct.Int", 1),
				rift.Path("Any.Int", 1),
				rift.Path("Map.Int", 1),
				rift.Path("Map.Arr.1", 1),
				rift.Path("Map.Arr.0", 2),
			},
			Then: &TestData{
				Int:      2,
				IntPtr:   ptr(3),
				String:   "Hi",
				Slice:    []TestData{{Int: 1}},
				SlicePtr: []*TestData{{Int: 1}},
				Struct:   &TestData{Int: 1},
				Any:      map[string]any{"Int": 1},
				Map:      map[string]any{"Int": 1, "Arr": []any{2, 1}},
			},
			Bnds: []rift.Change{
				{Path: "Int", Type: "int", New: 2, Old: 0},
				{Path: "IntPtr", Type: "int", New: 3, Old: nil},
				{Path: "String", Type: "string", New: "Hi", Old: ""},
				{Path: "Slice.0.Int", Type: "int", New: 1, Old: 0},
				{Path: "SlicePtr.0.Int", Type: "int", New: 1, Old: nil},
				{Path: "Struct.Int", Type: "int", New: 1, Old: nil},
				{Path: "Any.Int", Type: "int", New: 1, Old: nil},
				{Path: "Map.Int", Type: "int", New: 1, Old: nil},
				{Path: "Map.Arr.1", Type: "int", New: 1, Old: nil},
				{Path: "Map.Arr.0", Type: "int", New: 2, Old: nil},
			},
		},
		{
			Desc: "same as above, but tests Old Values.",
			Give: &TestData{
				Int:      11,
				IntPtr:   ptr(22),
				String:   "Hello",
				Slice:    []TestData{{Int: 33}},
				SlicePtr: []*TestData{{Int: 44}},
				Struct:   &TestData{Int: 55},
				Any:      map[string]any{"Int": 66},
				Map:      map[string]any{"Int": 77, "Arr": []any{88, 99}},
			},
			When: []rift.Node{
				rift.Path("Int", 2),
				rift.Path("IntPtr", 3),
				rift.Path("String", "Hi"),
				rift.Path("Slice.0.Int", 1),
				rift.Path("SlicePtr.0.Int", 1),
				rift.Path("Struct.Int", 1),
				rift.Path("Any.Int", 1),
				rift.Path("Map.Int", 1),
				rift.Path("Map.Arr.1", 1),
				rift.Path("Map.Arr.0", 2),
			},
			Then: &TestData{
				Int:      2,
				IntPtr:   ptr(3),
				String:   "Hi",
				Slice:    []TestData{{Int: 1}},
				SlicePtr: []*TestData{{Int: 1}},
				Struct:   &TestData{Int: 1},
				Any:      map[string]any{"Int": 1},
				Map:      map[string]any{"Int": 1, "Arr": []any{2, 1}},
			},
			Bnds: []rift.Change{
				{Path: "Int", Type: "int", New: 2, Old: 11},
				{Path: "IntPtr", Type: "int", New: 3, Old: 22},
				{Path: "String", Type: "string", New: "Hi", Old: "Hello"},
				{Path: "Slice.0.Int", Type: "int", New: 1, Old: 33},
				{Path: "SlicePtr.0.Int", Type: "int", New: 1, Old: 44},
				{Path: "Struct.Int", Type: "int", New: 1, Old: 55},
				{Path: "Any.Int", Type: "int", New: 1, Old: 66},
				{Path: "Map.Int", Type: "int", New: 1, Old: 77},
				{Path: "Map.Arr.1", Type: "int", New: 1, Old: 99},
				{Path: "Map.Arr.0", Type: "int", New: 2, Old: 88},
			},
		},
	}

	for _, tc := range tt {
		bs := rift.SetMany(&tc.Give, tc.When...)
		assertEqual(t, tc.Then, tc.Give, tc.Desc)
		assertEqual(t, tc.Bnds, bs, tc.Desc)
	}
}

func TestGetFlat(t *testing.T) {

	tt := []struct {
		Desc string
		Give any
		Then []rift.Node
	}{
		{
			Desc: "empty struct",
			Give: TestData{},
			Then: []rift.Node{
				{Path: "Int", Type: "int", Data: 0},
				{Path: "IntPtr", Type: "int", Data: nil},
				{Path: "String", Type: "string", Data: ""},
				{Path: "Slice", Type: "slice", Data: nil},
				{Path: "SlicePtr", Type: "slice", Data: nil},
				{Path: "Struct", Type: "struct", Data: nil},
				{Path: "Any", Type: "interface", Data: nil},
				{Path: "Map", Type: "map", Data: nil},
			},
		},
		{
			Desc: "filled struct",
			Give: TestData{
				Int:      11,
				IntPtr:   ptr(22),
				String:   "Hello",
				Slice:    []TestData{{Int: 33}},
				SlicePtr: []*TestData{{Int: 44}},
				Struct:   &TestData{Int: 55},
				Any:      map[string]any{"Int": 66},
				Map:      map[string]any{"Arr": []any{77, 88}},
			},
			Then: []rift.Node{
				{Path: "Int", Type: "int", Data: 11},
				{Path: "IntPtr", Type: "int", Data: 22},
				{Path: "String", Type: "string", Data: "Hello"},
				{Path: "Slice.0.Int", Type: "int", Data: 33},
				{Path: "Slice.0.IntPtr", Type: "int", Data: nil},
				{Path: "Slice.0.String", Type: "string", Data: ""},
				{Path: "Slice.0.Slice", Type: "slice", Data: nil},
				{Path: "Slice.0.SlicePtr", Type: "slice", Data: nil},
				{Path: "Slice.0.Struct", Type: "struct", Data: nil},
				{Path: "Slice.0.Any", Type: "interface", Data: nil},
				{Path: "Slice.0.Map", Type: "map", Data: nil},
				{Path: "SlicePtr.0.Int", Type: "int", Data: 44},
				{Path: "SlicePtr.0.IntPtr", Type: "int", Data: nil},
				{Path: "SlicePtr.0.String", Type: "string", Data: ""},
				{Path: "SlicePtr.0.Slice", Type: "slice", Data: nil},
				{Path: "SlicePtr.0.SlicePtr", Type: "slice", Data: nil},
				{Path: "SlicePtr.0.Struct", Type: "struct", Data: nil},
				{Path: "SlicePtr.0.Any", Type: "interface", Data: nil},
				{Path: "SlicePtr.0.Map", Type: "map", Data: nil},
				{Path: "Struct.Int", Type: "int", Data: 55},
				{Path: "Struct.IntPtr", Type: "int", Data: nil},
				{Path: "Struct.String", Type: "string", Data: ""},
				{Path: "Struct.Slice", Type: "slice", Data: nil},
				{Path: "Struct.SlicePtr", Type: "slice", Data: nil},
				{Path: "Struct.Struct", Type: "struct", Data: nil},
				{Path: "Struct.Any", Type: "interface", Data: nil},
				{Path: "Struct.Map", Type: "map", Data: nil},
				{Path: "Any.Int", Type: "int", Data: 66},
				{Path: "Map.Arr.0", Type: "int", Data: 77},
				{Path: "Map.Arr.1", Type: "int", Data: 88},
			},
		},
	}

	for _, tc := range tt {
		bs := rift.GetFlat(&tc.Give)
		assertEqual(t, tc.Then, bs, tc.Desc)
	}
}

func TestGet(t *testing.T) {

	tt := []struct {
		Desc string
		Give any
		Then any
	}{
		{
			Desc: "nil root",
			Give: nil,
			Then: rift.Node{Type: "interface"},
		},
		{
			Desc: "int root",
			Give: 3,
			Then: rift.Node{Type: "int", Data: 3},
		},
		{
			Desc: "zero struct",
			Give: TestData{},
			Then: rift.Node{
				Type: "struct",
				Next: []rift.Node{
					{Name: "Int", Path: "Int", Type: "int", Data: 0},
					{Name: "IntPtr", Path: "IntPtr", Type: "int", Data: nil},
					{Name: "String", Path: "String", Type: "string", Data: ""},
					{Name: "Slice", Path: "Slice", Type: "slice", Data: nil},
					{Name: "SlicePtr", Path: "SlicePtr", Type: "slice", Data: nil},
					{Name: "Struct", Path: "Struct", Type: "struct", Data: nil},
					{Name: "Any", Path: "Any", Type: "interface", Data: nil},
					{Name: "Map", Path: "Map", Type: "map", Data: nil},
				},
			},
		},
		{
			Desc: "filled struct",
			Give: TestData{
				Int:      11,
				IntPtr:   ptr(22),
				String:   "Hello",
				Slice:    []TestData{{Int: 33}},
				SlicePtr: []*TestData{{Int: 44}},
				Struct:   &TestData{Int: 55},
				Any:      map[string]any{"Int": 66},
				Map:      map[string]any{"Arr": []any{77, 88}},
			},
			Then: rift.Node{
				Type: "struct",
				Next: []rift.Node{
					{Name: "Int", Path: "Int", Type: "int", Data: 11},
					{Name: "IntPtr", Path: "IntPtr", Type: "int", Data: 22},
					{Name: "String", Path: "String", Type: "string", Data: "Hello"},
					{Name: "Slice", Path: "Slice", Type: "slice", Data: nil, Next: []rift.Node{
						{Name: "0", Path: "Slice.0", Type: "struct", Data: nil, Next: []rift.Node{
							{Name: "Int", Path: "Slice.0.Int", Type: "int", Data: 33},
							{Name: "IntPtr", Path: "Slice.0.IntPtr", Type: "int", Data: nil},
							{Name: "String", Path: "Slice.0.String", Type: "string", Data: ""},
							{Name: "Slice", Path: "Slice.0.Slice", Type: "slice", Data: nil},
							{Name: "SlicePtr", Path: "Slice.0.SlicePtr", Type: "slice", Data: nil},
							{Name: "Struct", Path: "Slice.0.Struct", Type: "struct", Data: nil},
							{Name: "Any", Path: "Slice.0.Any", Type: "interface", Data: nil},
							{Name: "Map", Path: "Slice.0.Map", Type: "map", Data: nil},
						}},
					}},
					{Name: "SlicePtr", Path: "SlicePtr", Type: "slice", Data: nil, Next: []rift.Node{
						{Name: "0", Path: "SlicePtr.0", Type: "struct", Data: nil, Next: []rift.Node{
							{Name: "Int", Path: "SlicePtr.0.Int", Type: "int", Data: 44},
							{Name: "IntPtr", Path: "SlicePtr.0.IntPtr", Type: "int", Data: nil},
							{Name: "String", Path: "SlicePtr.0.String", Type: "string", Data: ""},
							{Name: "Slice", Path: "SlicePtr.0.Slice", Type: "slice", Data: nil},
							{Name: "SlicePtr", Path: "SlicePtr.0.SlicePtr", Type: "slice", Data: nil},
							{Name: "Struct", Path: "SlicePtr.0.Struct", Type: "struct", Data: nil},
							{Name: "Any", Path: "SlicePtr.0.Any", Type: "interface", Data: nil},
							{Name: "Map", Path: "SlicePtr.0.Map", Type: "map", Data: nil},
						}},
					}},
					{Name: "Struct", Path: "Struct", Type: "struct", Data: nil, Next: []rift.Node{
						{Name: "Int", Path: "Struct.Int", Type: "int", Data: 55},
						{Name: "IntPtr", Path: "Struct.IntPtr", Type: "int", Data: nil},
						{Name: "String", Path: "Struct.String", Type: "string", Data: ""},
						{Name: "Slice", Path: "Struct.Slice", Type: "slice", Data: nil},
						{Name: "SlicePtr", Path: "Struct.SlicePtr", Type: "slice", Data: nil},
						{Name: "Struct", Path: "Struct.Struct", Type: "struct", Data: nil},
						{Name: "Any", Path: "Struct.Any", Type: "interface", Data: nil},
						{Name: "Map", Path: "Struct.Map", Type: "map", Data: nil},
					}},
					{Name: "Any", Path: "Any", Type: "map", Data: nil, Next: []rift.Node{
						{Name: "Int", Path: "Any.Int", Type: "int", Data: 66},
					}},
					{Name: "Map", Path: "Map", Type: "map", Data: nil, Next: []rift.Node{
						{Name: "Arr", Path: "Map.Arr", Type: "slice", Data: nil, Next: []rift.Node{
							{Name: "0", Path: "Map.Arr.0", Type: "int", Data: 77},
							{Name: "1", Path: "Map.Arr.1", Type: "int", Data: 88},
						}},
					}},
				},
			},
		},
	}

	for _, tc := range tt {
		tree := rift.Get(tc.Give)
		assertEqual(t, tc.Then, tree, tc.Desc)
	}
}

func TestSet(t *testing.T) {

	tt := []struct {
		Desc string
		Give any
		When rift.Node
		Then any
		Chng []rift.Change
	}{
		{
			Desc: "set a field",
			Give: &TestData{},
			When: rift.Node{Path: "Int", Data: 3},
			Then: &TestData{Int: 3},
			Chng: []rift.Change{
				{Path: "Int", Type: "int", New: 3, Old: 0},
			},
		},
		{
			Desc: "set subfields",
			Give: &TestData{Struct: &TestData{Int: 3, String: "A"}},
			When: rift.Node{Next: []rift.Node{{Path: "Struct.Int", Data: 4}, {Path: "Struct.String", Data: "B"}}},
			Then: &TestData{Struct: &TestData{Int: 4, String: "B"}},
			Chng: []rift.Change{
				{Path: "Struct.Int", Type: "int", New: 4, Old: 3},
				{Path: "Struct.String", Type: "string", New: "B", Old: "A"},
			},
		},
	}

	for _, tc := range tt {
		cs := rift.Set(tc.Give, tc.When)
		assertEqual(t, tc.Then, tc.Give, tc.Desc)
		assertEqual(t, tc.Chng, cs, tc.Desc)
	}
}

type TestData struct {
	Int      int
	IntPtr   *int
	String   string
	Slice    []TestData
	SlicePtr []*TestData
	Struct   *TestData
	Any      any
	Map      map[string]any
}

func assertEqual(t *testing.T, exp, got any, msgs ...any) {
	t.Helper()
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("\nExp:\n%v\nGot:\n%v\nMsg: %v", exp, got, fmt.Sprint(msgs...))
	}
}

func ptr[T any](v T) *T {
	return &v
}
