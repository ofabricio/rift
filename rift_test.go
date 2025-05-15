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
		fmt.Println(v.Path, v.Type, v.OldValue, v.NewValue)
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
		When []rift.F
		Then TestData
	}{
		{
			When: []rift.F{
				rift.Field("Int", 3),
				rift.Field("String", "Test"),
				rift.Field("Slice.1.Int", 2),
				rift.Field("Slice.0.Int", 1),
				rift.Field("Slice.2.Int", 3),
				rift.Field("SlicePtr.1.Int", 5),
				rift.Field("SlicePtr.0.Int", 4),
				rift.Field("SlicePtr.2.Int", 6),
				rift.Field("Struct.Slice.1.Int", 123),
			},
			Then: TestData{
				Int:      3,
				String:   "Test",
				Slice:    []TestData{{Int: 1}, {Int: 2}, {Int: 3}},
				SlicePtr: []*TestData{{Int: 4}, {Int: 5}, {Int: 6}},
				Struct:   &TestData{Slice: []TestData{{}, {Int: 123}}},
			},
		},
		{
			When: []rift.F{
				rift.Field("Slice", []TestData{{Int: 1}, {Int: 2}}),
				rift.Field("SlicePtr", []*TestData{{Int: 3}, {Int: 4}}),
				rift.Field("Struct", &TestData{Int: 5, Slice: []TestData{{Int: 6}}}),
			},
			Then: TestData{
				Slice:    []TestData{{Int: 1}, {Int: 2}},
				SlicePtr: []*TestData{{Int: 3}, {Int: 4}},
				Struct:   &TestData{Int: 5, Slice: []TestData{{Int: 6}}},
			},
		},
	}

	for i, tc := range tt {
		got := TestData{}
		rift.Bind(&got, tc.When...)
		assertEqual(t, tc.Then, got, i)
	}
}

type TestData struct {
	Int      int
	String   string
	Slice    []TestData
	SlicePtr []*TestData
	Struct   *TestData
}

func assertEqual(t *testing.T, exp, got any, msgs ...any) {
	t.Helper()
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("\nExp:\n%v\nGot:\n%v\nMsg: %v", exp, got, fmt.Sprint(msgs...))
	}
}
