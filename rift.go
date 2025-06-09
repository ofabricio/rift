package rift

import (
	"reflect"
	"strconv"
	"strings"
)

// Unbind extracts public fields from a struct.
// It returns a slice of unbound values that
// contain the field's path, type, and value.
func Unbind(v any) []Unbound {
	var out []Unbound
	unbind(reflect.ValueOf(v), "", &out)
	return out
}

func unbind(v reflect.Value, path string, out *[]Unbound) {
	switch v.Kind() {
	case reflect.Invalid:
		*out = append(*out, Field(path, nil))
	case reflect.Interface:
		unbind(v.Elem(), path, out)
	case reflect.Pointer:
		if v.IsNil() {
			*out = append(*out, fieldWithType(path, nil, v.Type().Elem().Kind().String()))
			return
		}
		unbind(v.Elem(), path, out)
	case reflect.Struct:
		for i := range v.NumField() {
			f := v.Field(i)
			p := v.Type().Field(i).Name
			unbind(f, joinPath(path, p), out)
		}
	case reflect.Map:
		if v.Len() == 0 {
			*out = append(*out, fieldWithType(path, nil, reflect.Map.String()))
			return
		}
		for iter := v.MapRange(); iter.Next(); {
			k := iter.Key()
			v := iter.Value()
			p := k.String()
			unbind(v, joinPath(path, p), out)
		}
	case reflect.Slice:
		if v.Len() == 0 {
			*out = append(*out, fieldWithType(path, nil, reflect.Slice.String()))
			return
		}
		for i := range v.Len() {
			f := v.Index(i)
			p := strconv.Itoa(i)
			unbind(f, joinPath(path, p), out)
		}
	default:
		*out = append(*out, Field(path, v.Interface()))
	}
}

func joinPath(path, subpath string) string {
	if path != "" {
		return path + "." + subpath
	}
	return subpath
}

// Bind sets values to a struct based on the provided fields.
// It returns a slice of bound values, that contain the path,
// the type, the new value set, and the old value before set.
func Bind(dst any, fs ...Unbound) []Bound {
	bds := make([]Bound, 0, len(fs))
	vOf := reflect.ValueOf(dst)
	for _, f := range fs {
		old := bind(vOf, reflect.ValueOf(f.Data), f.Path)
		bnd := Bound{Path: f.Path, Type: f.Type, New: f.Data, Old: old}
		bds = append(bds, bnd)
	}
	return bds
}

func bind(dst, val reflect.Value, path string) (old any) {

	keyOrIdx, rest, _ := strings.Cut(path, ".")

	switch dst.Kind() {
	case reflect.Invalid:
	case reflect.Pointer:
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
			_ = bind(dst.Elem(), val, path)
			old = nil
		} else {
			old = bind(dst.Elem(), val, path)
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
			old = bind(dst.Elem().Index(n), val, rest)
		} else {
			if dst.IsNil() {
				new := reflect.MakeMap(reflect.TypeFor[map[string]any]())
				dst.Set(new)
			}
			old = bind(dst.Elem(), val, path)
		}
	case reflect.Slice:
		n, _ := getNumber(keyOrIdx)
		if n >= dst.Len() {
			new := reflect.MakeSlice(dst.Type(), n+1, n+1)
			reflect.Copy(new, dst)
			dst.Set(new)
		}
		old = bind(dst.Index(n), val, rest)
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
		old = bind(new, val, rest)
		dst.SetMapIndex(k, new.Elem())
	case reflect.Struct:
		key := dst.FieldByName(keyOrIdx)
		old = bind(key, val, rest)
	default:
		old = dst.Interface()
		dst.Set(val)
	}
	return
}

// Field creates a field with the specified path and value.
func Field(path string, value any) Unbound {
	if v := reflect.ValueOf(value); v.IsValid() {
		return fieldWithType(path, value, v.Kind().String())
	}
	return fieldWithType(path, value, reflect.Interface.String())
}

// Field creates a field with the specified path and value.
func fieldWithType(path string, value any, typ string) Unbound {
	return Unbound{
		Path: path,
		Type: typ,
		Data: value,
	}
}

func getNumber(path string) (int, bool) {
	v, err := strconv.Atoi(path)
	return v, err == nil
}

// Unbound represents a field that is not yet bound to a struct.
type Unbound struct {
	Path string
	Type string
	Data any
}

// Bound represents a field that has been bound to a struct.
type Bound struct {
	Path string
	Type string
	New  any // New value set.
	Old  any // Old value before set.
}

// Describe returns a tree representation of the provided value.
func Describe(v any) Tree {
	var t Tree
	describe(reflect.ValueOf(v), "", &t)
	return t
}

func describe(v reflect.Value, path string, out *Tree) {
	out.Path = path
	out.Type = v.Kind().String()
	switch v.Kind() {
	case reflect.Invalid:
		out.Type = reflect.Interface.String()
	case reflect.Interface:
		describe(v.Elem(), path, out)
	case reflect.Pointer:
		if v.IsNil() {
			out.Type = v.Type().Elem().Kind().String()
			return
		}
		describe(v.Elem(), path, out)
	case reflect.Slice:
		for i := range v.Len() {
			f := v.Index(i)
			p := strconv.Itoa(i)
			n := Tree{Name: p}
			describe(f, joinPath(path, p), &n)
			out.Next = append(out.Next, n)
		}
	case reflect.Map:
		for iter := v.MapRange(); iter.Next(); {
			k := iter.Key()
			v := iter.Value()
			p := k.String()
			n := Tree{Name: p}
			describe(v, joinPath(path, p), &n)
			out.Next = append(out.Next, n)
		}
	case reflect.Struct:
		for i := range v.NumField() {
			f := v.Field(i)
			p := v.Type().Field(i).Name
			n := Tree{Name: p}
			describe(f, joinPath(path, p), &n)
			out.Next = append(out.Next, n)
		}
	default:
		out.Data = v.Interface()
	}
}

type Tree struct {
	Name string
	Path string
	Type string
	Data any
	Next []Tree
}
