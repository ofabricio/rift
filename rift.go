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
	return unbind(reflect.ValueOf(v), "")
}

func unbind(v reflect.Value, path string) (vs []Unbound) {
	switch v.Kind() {
	case reflect.Pointer:
		vs = append(vs, unbind(reflect.Indirect(v), path)...)
	case reflect.Struct:
		for i := range v.NumField() {
			f := v.Field(i)
			p := v.Type().Field(i).Name
			vs = append(vs, unbind(f, joinPath(path, p))...)
		}
	case reflect.Slice:
		for i := range v.Len() {
			f := v.Index(i)
			p := strconv.Itoa(i)
			vs = append(vs, unbind(f, joinPath(path, p))...)
		}
	default:
		vs = append(vs, Field(path, v.Interface()))
	}
	return
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
		old := bind(vOf, reflect.ValueOf(f.Value), f.Path)
		bnd := Bound{Path: f.Path, Type: f.Type, New: f.Value, Old: old}
		bds = append(bds, bnd)
	}
	return bds
}

func bind(dst, val reflect.Value, path string) (old any) {

	ptr := dst
	dst = reflect.Indirect(dst)

	keyOrIdx, rest, _ := strings.Cut(path, ".")

	switch dst.Kind() {
	case reflect.Interface:
		if path == "" {
			old = dst.Interface()
			dst.Set(val)
			return
		}
		if n, ok := getNumber(keyOrIdx); ok {
			if dst.IsNil() {
				new := reflect.MakeSlice(reflect.TypeFor[[]any](), n+1, n+1)
				dst.Set(new)
				old = bind(dst.Elem().Index(n), val, rest)
			} else if n >= dst.Elem().Len() {
				new := reflect.MakeSlice(dst.Elem().Type(), n+1, n+1)
				reflect.Copy(new, dst.Elem())
				dst.Set(new)
				old = bind(new.Index(n), val, rest)
			} else {
				old = bind(dst.Elem().Index(n), val, rest)
			}
		} else {
			if dst.IsNil() {
				new := reflect.MakeMap(reflect.TypeFor[map[string]any]())
				dst.Set(new)
				old = bind(new, val, path)
			} else {
				old = bind(dst.Elem(), val, path)
			}
		}
		return
	case reflect.Slice:
		n, _ := getNumber(keyOrIdx)
		if n < dst.Len() {
			old = bind(dst.Index(n), val, rest)
		} else {
			newSlice := reflect.MakeSlice(dst.Type(), n+1, n+1)
			reflect.Copy(newSlice, dst)
			dst.Set(newSlice)
			old = bind(dst.Index(n), val, rest)
		}
		return
	case reflect.Map:
		if dst.IsNil() {
			new := reflect.MakeMap(dst.Type())
			dst.Set(new)
		}
		key := reflect.ValueOf(keyOrIdx)
		item := dst.MapIndex(key)
		if !item.IsValid() {
			item = reflect.New(dst.Type().Elem()).Elem()
		}
		n := reflect.New(item.Type()).Elem()
		n.Set(item)
		old = bind(n, val, rest)
		dst.SetMapIndex(key, n.Elem())
		return
	case reflect.Struct:
		field := dst.FieldByName(keyOrIdx)
		old = bind(field, val, rest)
		return
	}
	if ptr.Kind() == reflect.Pointer {
		if ptr.IsNil() {
			n := reflect.New(ptr.Type().Elem())
			_ = bind(n.Elem(), val, path)
			ptr.Set(n)
			old = nil
		} else {
			old = bind(ptr.Elem(), val, path)
		}
	} else {
		old = dst.Interface()
		dst.Set(val)
	}
	return
}

// Field creates a field with the specified path and value.
func Field(path string, value any) Unbound {
	return Unbound{
		Path:  path,
		Type:  reflect.TypeOf(value).Name(),
		Value: value,
	}
}

// Unbound represents a field that is not yet bound to a struct.
type Unbound struct {
	Path  string
	Type  string
	Value any
}

// Bound represents a field that has been bound to a struct.
type Bound struct {
	Path string
	Type string
	New  any // New value set.
	Old  any // Old value before set.
}

func getNumber(path string) (int, bool) {
	v, err := strconv.Atoi(path)
	return v, err == nil
}
