package rift_test

import (
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

	bs := rift.Bind(&user,
		rift.Field("Name", "John"),
		rift.Field("Addresses.0.Street", "Main"),
		rift.Field("Addresses.0.Number", 100),
		rift.Field("Addresses.1.Street", "Avenue"),
		rift.Field("Addresses.1.Number", 200),
	)

	fmt.Println("User:")
	fmt.Println(user)

	fmt.Println("Bound:")
	for _, v := range bs {
		fmt.Println(v.Path, v.Type, v.Old, v.New)
	}

	fmt.Println("Unbind:")
	for _, v := range rift.Unbind(&user) {
		fmt.Println(v.Path, v.Type, v.Value)
	}

	// Output:
	// User:
	// {John [{Main 100} {Avenue 200}]}
	// Bound:
	// Name string Luke John
	// Addresses.0.Street string  Main
	// Addresses.0.Number int 0 100
	// Addresses.1.Street string  Avenue
	// Addresses.1.Number int 0 200
	// Unbind:
	// Name string John
	// Addresses.0.Street string Main
	// Addresses.0.Number int 100
	// Addresses.1.Street string Avenue
	// Addresses.1.Number int 200
}

func TestBind(t *testing.T) {

	tt := []struct {
		Desc string
		Give any
		When []rift.Unbound
		Then any
		Bnds []rift.Bound
	}{
		{
			Desc: "given a nil source, it should set the root if field is empty",
			Give: nil,
			When: []rift.Unbound{
				rift.Field("", 3),
			},
			Then: 3,
			Bnds: []rift.Bound{{Path: "", Type: "int", New: 3, Old: nil}},
		},
		{
			Desc: "given a nil source, it should set the root as a slice if field is number",
			Give: nil,
			When: []rift.Unbound{
				rift.Field("0", 3),
			},
			Then: []any{3},
			Bnds: []rift.Bound{{Path: "0", Type: "int", New: 3, Old: nil}},
		},
		{
			Desc: "given a nil source, it should set the root as a slice of len=2 if field is number 1",
			Give: nil,
			When: []rift.Unbound{
				rift.Field("1", 3),
			},
			Then: []any{nil, 3},
			Bnds: []rift.Bound{{Path: "1", Type: "int", New: 3, Old: nil}},
		},
		{
			Desc: "given a nil source, it should set the root as a slice of len=2",
			Give: nil,
			When: []rift.Unbound{
				rift.Field("0", 2),
				rift.Field("1", 3),
			},
			Then: []any{2, 3},
			Bnds: []rift.Bound{
				{Path: "0", Type: "int", New: 2, Old: nil},
				{Path: "1", Type: "int", New: 3, Old: nil},
			},
		},
		{
			Desc: "test inverted slice indexes",
			Give: nil,
			When: []rift.Unbound{
				rift.Field("1", 2),
				rift.Field("0", 3),
			},
			Then: []any{3, 2},
			Bnds: []rift.Bound{
				{Path: "1", Type: "int", New: 2, Old: nil},
				{Path: "0", Type: "int", New: 3, Old: nil},
			},
		},
		{
			Desc: "given a nil source, it should set the root as a map with key a",
			Give: nil,
			When: []rift.Unbound{
				rift.Field("a", 2),
			},
			Then: map[string]any{"a": 2},
			Bnds: []rift.Bound{
				{Path: "a", Type: "int", New: 2, Old: nil},
			},
		},
		{
			Desc: "given a nil source, it should set the root as a map with key a and b",
			Give: nil,
			When: []rift.Unbound{
				rift.Field("a", 2),
				rift.Field("b", 3),
			},
			Then: map[string]any{"a": 2, "b": 3},
			Bnds: []rift.Bound{
				{Path: "a", Type: "int", New: 2, Old: nil},
				{Path: "b", Type: "int", New: 3, Old: nil},
			},
		},
		{
			Desc: "given a nil source, it should set the root as a map with key a and b; and key c should also be a map",
			Give: nil,
			When: []rift.Unbound{
				rift.Field("a", 2),
				rift.Field("b", 3),
				rift.Field("c.a", 3),
			},
			Then: map[string]any{"a": 2, "b": 3, "c": map[string]any{"a": 3}},
			Bnds: []rift.Bound{
				{Path: "a", Type: "int", New: 2, Old: nil},
				{Path: "b", Type: "int", New: 3, Old: nil},
				{Path: "c.a", Type: "int", New: 3, Old: nil},
			},
		},
		{
			Desc: "given a nil source, it should set the root as a map with key a and b; and key c should also be a map with key a and b",
			Give: nil,
			When: []rift.Unbound{
				rift.Field("a", 2),
				rift.Field("b", 3),
				rift.Field("c.a", 4),
				rift.Field("c.b", 5),
			},
			Then: map[string]any{"a": 2, "b": 3, "c": map[string]any{"a": 4, "b": 5}},
			Bnds: []rift.Bound{
				{Path: "a", Type: "int", New: 2, Old: nil},
				{Path: "b", Type: "int", New: 3, Old: nil},
				{Path: "c.a", Type: "int", New: 4, Old: nil},
				{Path: "c.b", Type: "int", New: 5, Old: nil},
			},
		},
		{
			Desc: "if a subfield of a generic map key has an index, that field should be a slice",
			Give: nil,
			When: []rift.Unbound{
				rift.Field("a.b.0", 3),
			},
			Then: map[string]any{"a": map[string]any{"b": []any{3}}},
			Bnds: []rift.Bound{
				{Path: "a.b.0", Type: "int", New: 3, Old: nil},
			},
		},
		{
			Desc: "same as above, but with index 1",
			Give: nil,
			When: []rift.Unbound{
				rift.Field("a.b.1", 3),
			},
			Then: map[string]any{"a": map[string]any{"b": []any{nil, 3}}},
			Bnds: []rift.Bound{
				{Path: "a.b.1", Type: "int", New: 3, Old: nil},
			},
		},
		{
			Desc: "same as above, but with index 0 and 1",
			Give: nil,
			When: []rift.Unbound{
				rift.Field("a.a.0", 2),
				rift.Field("a.a.1", 3),
			},
			Then: map[string]any{"a": map[string]any{"a": []any{2, 3}}},
			Bnds: []rift.Bound{
				{Path: "a.a.0", Type: "int", New: 2, Old: nil},
				{Path: "a.a.1", Type: "int", New: 3, Old: nil},
			},
		},
		{
			Desc: "same as above, but with inverted indexes",
			Give: nil,
			When: []rift.Unbound{
				rift.Field("a.a.1", 2),
				rift.Field("a.a.0", 3),
			},
			Then: map[string]any{"a": map[string]any{"a": []any{3, 2}}},
			Bnds: []rift.Bound{
				{Path: "a.a.1", Type: "int", New: 2, Old: nil},
				{Path: "a.a.0", Type: "int", New: 3, Old: nil},
			},
		},
		{
			Desc: "now a slice with a map inside",
			Give: nil,
			When: []rift.Unbound{
				rift.Field("a.a.0", 2),
				rift.Field("a.a.1.a", 3),
				rift.Field("a.a.1.b", 4),
			},
			Then: map[string]any{"a": map[string]any{"a": []any{2, map[string]any{"a": 3, "b": 4}}}},
			Bnds: []rift.Bound{
				{Path: "a.a.0", Type: "int", New: 2, Old: nil},
				{Path: "a.a.1.a", Type: "int", New: 3, Old: nil},
				{Path: "a.a.1.b", Type: "int", New: 4, Old: nil},
			},
		},
		{
			Desc: "test setting many different field types of a struct",
			Give: &TestData{},
			When: []rift.Unbound{
				rift.Field("Int", 2),
				rift.Field("IntPtr", 3),
				rift.Field("String", "Hi"),
				rift.Field("Slice.0.Int", 1),
				rift.Field("SlicePtr.0.Int", 1),
				rift.Field("Struct.Int", 1),
				rift.Field("Any.Int", 1),
				rift.Field("Map.Int", 1),
				rift.Field("Map.Arr.1", 1),
				rift.Field("Map.Arr.0", 2),
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
			Bnds: []rift.Bound{
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
			When: []rift.Unbound{
				rift.Field("Int", 2),
				rift.Field("IntPtr", 3),
				rift.Field("String", "Hi"),
				rift.Field("Slice.0.Int", 1),
				rift.Field("SlicePtr.0.Int", 1),
				rift.Field("Struct.Int", 1),
				rift.Field("Any.Int", 1),
				rift.Field("Map.Int", 1),
				rift.Field("Map.Arr.1", 1),
				rift.Field("Map.Arr.0", 2),
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
			Bnds: []rift.Bound{
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
		bs := rift.Bind(&tc.Give, tc.When...)
		assertEqual(t, tc.Then, tc.Give, tc.Desc)
		assertEqual(t, tc.Bnds, bs, tc.Desc)
	}
}

func TestUnbind(t *testing.T) {

	tt := []struct {
		Desc string
		Give any
		Then []rift.Unbound
	}{
		{
			Desc: "Unbind empty struct",
			Give: TestData{},
			Then: []rift.Unbound{
				{Path: "Int", Type: "int", Value: 0},
				{Path: "IntPtr", Type: "nil", Value: nil},
				{Path: "String", Type: "string", Value: ""},
				{Path: "Struct", Type: "nil", Value: nil},
				{Path: "Any", Type: "nil", Value: nil},
			},
		},
		{
			Desc: "Unbind filled struct",
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
			Then: []rift.Unbound{
				{Path: "Int", Type: "int", Value: 11},
				{Path: "IntPtr", Type: "int", Value: 22},
				{Path: "String", Type: "string", Value: "Hello"},
				{Path: "Slice.0.Int", Type: "int", Value: 33},
				{Path: "Slice.0.IntPtr", Type: "nil", Value: nil},
				{Path: "Slice.0.String", Type: "string", Value: ""},
				{Path: "Slice.0.Struct", Type: "nil", Value: nil},
				{Path: "Slice.0.Any", Type: "nil", Value: nil},
				{Path: "SlicePtr.0.Int", Type: "int", Value: 44},
				{Path: "SlicePtr.0.IntPtr", Type: "nil", Value: nil},
				{Path: "SlicePtr.0.String", Type: "string", Value: ""},
				{Path: "SlicePtr.0.Struct", Type: "nil", Value: nil},
				{Path: "SlicePtr.0.Any", Type: "nil", Value: nil},
				{Path: "Struct.Int", Type: "int", Value: 55},
				{Path: "Struct.IntPtr", Type: "nil", Value: nil},
				{Path: "Struct.String", Type: "string", Value: ""},
				{Path: "Struct.Struct", Type: "nil", Value: nil},
				{Path: "Struct.Any", Type: "nil", Value: nil},
				{Path: "Any.Int", Type: "int", Value: 66},
				{Path: "Map.Arr.0", Type: "int", Value: 77},
				{Path: "Map.Arr.1", Type: "int", Value: 88},
			},
		},
	}

	for _, tc := range tt {
		bs := rift.Unbind(&tc.Give)
		assertEqual(t, tc.Then, bs, tc.Desc)
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
