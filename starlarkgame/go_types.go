package starlarkgame

import (
	"crypto/sha256"
	"fmt"
	"hash/fnv"
	"reflect"

	"go.starlark.net/starlark"
)

// A Go struct type wrapped as a starlark value.
type goStruct struct {
	typename string
	v        reflect.Value
}

// Wraps a Go struct type as a starlark.Value.
//
// Only fields with a "starlark" tag are visible to starlark code.
// The tag's value is used as the getter name. For example, given
// the following struct type:
//
//     type Object struct {
//         Foo string `starlark:"foo"`
//         Bar string
//     }
//
// The first line of the following starlark code will print the value
// of the `Foo` property of the given object, but the second line will
// err because the attribute is not found:
//
//     print(object.foo)
//     print(object.bar)
func newGoStruct(typename string, object interface{}) *goStruct {
	v := reflect.ValueOf(object)
	if v.Kind() != reflect.Struct && !(v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct) {
		panic(fmt.Errorf("value must be a struct or pointer to a struct, but was %T", v.Interface()))
	}
	return &goStruct{
		typename: typename,
		v:        v,
	}
}

func (s *goStruct) Attr(name string) (starlark.Value, error) {
	for i, f := range s.fields() {
		if f != name {
			continue
		}
		g := getter(s.v, f, i)
		return g.BindReceiver(s), nil
	}
	return nil, nil
}

func (s *goStruct) AttrNames() []string {
	var names []string
	for _, f := range s.fields() {
		if f != "" {
			names = append(names, f)
		}
	}
	return names
}

func (s *goStruct) fields() []string {
	var fs []string

	t := s.v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("starlark")
		if tag == "" {
			fs = append(fs, "")
			continue
		}
		fs = append(fs, tag)
	}

	return fs
}

func (s *goStruct) Freeze() {
	return // TODO: Implement
}

func (s *goStruct) Hash() (uint32, error) {
	sum := sha256.New().Sum([]byte(s.String()))
	h := fnv.New32a()
	h.Write(sum)
	return h.Sum32(), nil
}

func (s *goStruct) String() string {
	return ""
}

func (s *goStruct) Type() string {
	return s.typename
}

func (s *goStruct) Truth() starlark.Bool {
	return starlark.True
}

func getter(receiver reflect.Value, fieldName string, fieldIndex int) *starlark.Builtin {
	builtin := func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		// Ensure no args were passed.
		if err := starlark.UnpackPositionalArgs(f.Name(), args, kwargs, 0); err != nil {
			return nil, err
		}

		v := f.Receiver().(*goStruct)
		field := v.v.Field(fieldIndex)
		wrapper := selectTypeWrapper(field.Type())

		return wrapper(field.Interface()), nil
	}

	return starlark.NewBuiltin(fieldName, builtin)
}

type typeWrapper func(interface{}) starlark.Value

func selectTypeWrapper(t reflect.Type) typeWrapper {
	switch t.Kind() {
	case reflect.String:
		return func(v interface{}) starlark.Value { return starlark.String(v.(string)) }
	case reflect.Int:
		return func(v interface{}) starlark.Value { return starlark.MakeInt(v.(int)) }
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(v interface{}) starlark.Value { return starlark.MakeInt64(v.(int64)) }
	case reflect.Float32, reflect.Float64:
		return func(v interface{}) starlark.Value { return starlark.Float(v.(float64)) }
	}
	return nil
}
