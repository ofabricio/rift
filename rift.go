package rift

import (
	"reflect"
	"strconv"
	"strings"
)

// Unbind extracts public fields from a struct.
// It returns a slice of values that contain the field's path, type, and value.
func Unbind(v any) []F {
	return unbind(reflect.ValueOf(v), "")
}

func unbind(v reflect.Value, path string) (fs []F) {
	switch v.Kind() {
	case reflect.Pointer:
		fs = append(fs, unbind(reflect.Indirect(v), path)...)
	case reflect.Struct:
		for i := range v.NumField() {
			f := v.Field(i)
			p := v.Type().Field(i).Name
			fs = append(fs, unbind(f, joinPath(path, p))...)
		}
	case reflect.Slice:
		for i := range v.Len() {
			f := v.Index(i)
			p := strconv.Itoa(i)
			fs = append(fs, unbind(f, joinPath(path, p))...)
		}
	default:
		fs = append(fs, Field(path, v.Interface()))
	}
	return
}

func joinPath(path, subpath string) string {
	if path != "" {
		return path + "." + subpath
	}
	return subpath
}

// Bind binds values to a struct based on the provided paths.
// It returns a slice of bound values that contain the field's path, type, new value, and old value.
func Bind(dst any, fs ...F) []Bound {
	bs := make([]Bound, 0, len(fs))
	vOf := reflect.ValueOf(dst)
	for _, f := range fs {
		b := bind(vOf, reflect.ValueOf(f.Value), f.Path)
		bs = append(bs, b)
	}
	return bs
}

func bind(dst, val reflect.Value, path string) Bound {
	var b Bound
	if path == "" {
		b = Bound{Type: val.Type().Name(), NewValue: val.Interface(), OldValue: dst.Interface()}
		dst.Set(val)
		return b
	}
	nameOrIndex, nextPath, _ := strings.Cut(path, ".")
	switch dst.Kind() {
	case reflect.Interface:
		switch {
		case dst.Interface() == nil:
			dst.Set(reflect.MakeMap(reflect.TypeFor[map[string]any]()))
		}
		b = bind(dst.Elem(), val, path)
	case reflect.Pointer:
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		b = bind(reflect.Indirect(dst), val, path)
	case reflect.Struct:
		field := dst.FieldByName(nameOrIndex)
		b = bind(field, val, nextPath)
	case reflect.Slice:
		index, _ := strconv.Atoi(nameOrIndex)
		if index < dst.Len() {
			item := dst.Index(index)
			b = bind(item, val, nextPath)
		} else {
			item := reflect.New(dst.Type().Elem()).Elem()
			b = bind(item, val, nextPath)
			newSlice := reflect.MakeSlice(dst.Type(), index+1, index+1)
			reflect.Copy(newSlice, dst)
			newSlice.Index(index).Set(item)
			dst.Set(newSlice)
		}
	case reflect.Map:
		if dst.IsNil() {
			dst.Set(reflect.MakeMap(dst.Type()))
		}
		var elem reflect.Value
		key := reflect.ValueOf(nameOrIndex)
		if v := dst.MapIndex(key); v.IsValid() {
			elem = v
		} else if nextPath == "" {
			// If the key does not exist, create a new element
			// with the type of the map's value.
			elem = reflect.New(dst.Type().Elem()).Elem()
		} else {
			// If there is more path, then the element is a map.
			elem = reflect.MakeMap(reflect.TypeFor[map[string]any]())
		}
		b = bind(elem, val, nextPath)
		dst.SetMapIndex(key, elem)
	}
	b.Path = path
	return b
}

// Field creates a field with the specified path and value.
func Field(path string, value any) F {
	return F{
		Path:  path,
		Type:  reflect.TypeOf(value).Name(),
		Value: value,
	}
}

type F struct {
	Path  string
	Type  string
	Value any
}

type Bound struct {
	Path     string
	Type     string
	NewValue any
	OldValue any
}
