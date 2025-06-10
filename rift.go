package rift

import (
	"reflect"
	"strconv"
	"strings"
)

// Get returns a tree representation of the provided value.
func Get(v any) Node {
	var out Node
	get(reflect.ValueOf(v), "", &out)
	return out
}

func get(v reflect.Value, path string, out *Node) {
	out.Path = path
	out.Type = v.Kind().String()
	switch v.Kind() {
	case reflect.Invalid:
		out.Type = reflect.Interface.String()
	case reflect.Interface:
		get(v.Elem(), path, out)
	case reflect.Pointer:
		if v.IsNil() {
			out.Type = v.Type().Elem().Kind().String()
			return
		}
		get(v.Elem(), path, out)
	case reflect.Slice:
		for i := range v.Len() {
			f := v.Index(i)
			p := strconv.Itoa(i)
			n := Node{Name: p}
			get(f, joinPath(path, p), &n)
			out.Next = append(out.Next, n)
		}
	case reflect.Map:
		for iter := v.MapRange(); iter.Next(); {
			k := iter.Key()
			v := iter.Value()
			p := k.String()
			n := Node{Name: p}
			get(v, joinPath(path, p), &n)
			out.Next = append(out.Next, n)
		}
	case reflect.Struct:
		for i := range v.NumField() {
			f := v.Field(i)
			p := v.Type().Field(i).Name
			n := Node{Name: p}
			get(f, joinPath(path, p), &n)
			out.Next = append(out.Next, n)
		}
	default:
		out.Data = v.Interface()
	}
}

// GetFlat returns a flat representation of the provided value.
func GetFlat(v any) []Node {
	var out []Node
	walk(Get(v), func(n Node) {
		if len(n.Next) == 0 {
			n.Name = ""
			out = append(out, n)
		}
	})
	return out
}

// Set sets values to a struct based on the provided node.
func Set(dst any, n Node) []Change {
	var chgs []Change
	walk(n, func(n Node) {
		if len(n.Next) == 0 {
			chgs = append(chgs, SetPath(dst, n.Path, n.Data))
		}
	})
	return chgs
}

func walk(n Node, fn func(Node)) {
	for _, v := range n.Next {
		walk(v, fn)
	}
	fn(n)
}

// SetMany sets values to a struct based on the provided nodes.
func SetMany(dst any, ns ...Node) []Change {
	chg := make([]Change, 0, len(ns))
	for _, n := range ns {
		chg = append(chg, SetPath(dst, n.Path, n.Data))
	}
	return chg
}

// SetPath sets a value to a struct based on the provided path.
func SetPath(dst any, path string, val any) Change {
	d := reflect.ValueOf(dst)
	v := reflect.ValueOf(val)
	old := setPath(d, v, path)
	return Change{Path: path, New: val, Old: old, Type: getType(v)}
}

func setPath(dst, val reflect.Value, path string) (old any) {

	keyOrIdx, rest, _ := strings.Cut(path, ".")

	switch dst.Kind() {
	case reflect.Invalid:
	case reflect.Pointer:
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
			_ = setPath(dst.Elem(), val, path)
			old = nil
		} else {
			old = setPath(dst.Elem(), val, path)
		}
	case reflect.Interface:
		if path == "" {
			old = dst.Interface()
			dst.Set(val)
		} else if n, ok := getNumber(keyOrIdx); ok {
			if dst.IsNil() {
				new := reflect.MakeSlice(reflect.TypeFor[[]any](), n+1, n+1)
				dst.Set(new)
			} else if n >= dst.Elem().Len() {
				new := reflect.MakeSlice(dst.Elem().Type(), n+1, n+1)
				reflect.Copy(new, dst.Elem())
				dst.Set(new)
			}
			old = setPath(dst.Elem().Index(n), val, rest)
		} else {
			if dst.IsNil() {
				new := reflect.MakeMap(reflect.TypeFor[map[string]any]())
				dst.Set(new)
			}
			old = setPath(dst.Elem(), val, path)
		}
	case reflect.Slice:
		n, _ := getNumber(keyOrIdx)
		if n >= dst.Len() {
			new := reflect.MakeSlice(dst.Type(), n+1, n+1)
			reflect.Copy(new, dst)
			dst.Set(new)
		}
		old = setPath(dst.Index(n), val, rest)
	case reflect.Map:
		if dst.IsNil() {
			dst.Set(reflect.MakeMap(dst.Type()))
		}
		k := reflect.ValueOf(keyOrIdx)
		v := dst.MapIndex(k)
		if !v.IsValid() {
			v = reflect.New(dst.Type().Elem()).Elem()
		}
		new := reflect.New(v.Type()).Elem()
		new.Set(v)
		old = setPath(new, val, rest)
		dst.SetMapIndex(k, new.Elem())
	case reflect.Struct:
		key := dst.FieldByName(keyOrIdx)
		old = setPath(key, val, rest)
	default:
		old = dst.Interface()
		dst.Set(val)
	}
	return
}

// Path creates a node with the specified path and value.
func Path(path string, value any) Node {
	return Node{Path: path, Data: value}
}

// Change represents a change.
type Change struct {
	Path string
	Type string
	New  any // New value set.
	Old  any // Old value before set.
}

// Node is a node in the tree.
type Node struct {
	Name string
	Path string
	Type string
	Data any
	Next []Node
}

func joinPath(path, subpath string) string {
	if path != "" {
		return path + "." + subpath
	}
	return subpath
}

func getNumber(path string) (int, bool) {
	v, err := strconv.Atoi(path)
	return v, err == nil
}

func getType(v reflect.Value) string {
	if v.IsValid() {
		return v.Type().Name()
	}
	return reflect.Interface.String()
}
