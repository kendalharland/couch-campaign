package starlarkgame

import (
	"crypto/sha256"
	"fmt"
	"hash/fnv"
	"reflect"
	"strconv"
	"strings"

	"go.starlark.net/starlark"
)

// A Go struct type wrapped as a starlark value.
type goStruct struct {
	typename string
	builtins map[string]*starlark.Builtin
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
	if !(v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct) {
		panic(fmt.Errorf("value must be a pointer to a struct, but was %T", v.Interface()))
	}
	builtins, err := makeBuiltins(v)
	if err != nil {
		panic(err)
	}
	return &goStruct{
		typename: typename,
		builtins: builtins,
		v:        v,
	}
}

func (s *goStruct) Interface() interface{} {
	return s.v.Interface()
}

func (s *goStruct) Attr(name string) (starlark.Value, error) {
	attr, ok := s.builtins[name]
	if !ok {
		return nil, fmt.Errorf("not found %s", name)
	}
	return attr.BindReceiver(s), nil
}

func (s *goStruct) AttrNames() []string {
	var names []string
	for name := range s.builtins {
		names = append(names, name)
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

func makeBuiltins(v reflect.Value) (map[string]*starlark.Builtin, error) {
	builtins := make(map[string]*starlark.Builtin)

	t := v.Type().Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("starlark")
		if tag == "" {
			continue
		}
		options := strings.Split(tag, ",")
		getterName := options[0]
		builtins[getterName] = getter(v, getterName, i)

		if len(options) > 1 && options[1] == "mutable" {
			goSetterName := fmt.Sprintf("Set%s", field.Name)
			if _, found := reflect.PtrTo(t).MethodByName(goSetterName); !found {
				return nil, fmt.Errorf("&%s.%s is mutable but %s has no method named %q", t.Name(), field.Name, t.Name(), goSetterName)
			}
			setterName := fmt.Sprintf("set_%s", getterName)
			builtins[setterName] = setter(v, setterName, goSetterName, field.Type.Kind())
		}
	}

	return builtins, nil
}

func getter(_ reflect.Value, fieldName string, fieldIndex int) *starlark.Builtin {
	builtin := func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		// Ensure no args were passed.
		if err := starlark.UnpackPositionalArgs(f.Name(), args, kwargs, 0); err != nil {
			return nil, err
		}

		v := f.Receiver().(*goStruct)
		field := v.v.Elem().Field(fieldIndex)
		return wrapAsStarlarkValue(field), nil
	}

	return starlark.NewBuiltin(fieldName, builtin)
}

func setter(_ reflect.Value, name, goName string, kind reflect.Kind) *starlark.Builtin {
	builtin := func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var value starlark.Value
		if err := starlark.UnpackPositionalArgs(f.Name(), args, kwargs, 1, &value); err != nil {
			return nil, err
		}

		v := f.Receiver().(*goStruct)
		m := v.v.MethodByName(goName)
		params := []reflect.Value{
			reflect.ValueOf(unwrapStarlarkValue(value, kind)),
		}
		m.Call(params)
		return starlark.None, nil
	}

	return starlark.NewBuiltin(name, builtin)
}

func unwrapStarlarkValue(v starlark.Value, kind reflect.Kind) interface{} {
	switch kind {
	case reflect.String:
		return string(v.(starlark.String))
	case reflect.Int8:
		return int8(unwrapStarlarkInt(v, 8))
	case reflect.Int16:
		return int16(unwrapStarlarkInt(v, 16))
	case reflect.Int32:
		return int32(unwrapStarlarkInt(v, 32))
	case reflect.Int64:
		return int64(unwrapStarlarkInt(v, 64))
	case reflect.Int:
		return int(unwrapStarlarkInt(v, 64))
	case reflect.Float32:
		return float32(unwrapStarlarkFloat(v, 32))
	case reflect.Float64:
		return float64(unwrapStarlarkFloat(v, 32))
	}
	return nil
}

func unwrapStarlarkInt(v starlark.Value, bitSize int) int64 {
	i, err := strconv.ParseInt(v.String(), 10, bitSize)
	if err != nil {
		panic(err)
	}
	return i
}

func unwrapStarlarkFloat(v starlark.Value, bitSize int) float64 {
	f, err := strconv.ParseFloat(v.String(), bitSize)
	if err != nil {
		panic(err)
	}
	return f
}

func wrapAsStarlarkValue(f reflect.Value) starlark.Value {
	switch f.Type().Kind() {
	case reflect.String:
		v := f.Interface()
		return starlark.String(v.(string))
	case reflect.Int:
		v := f.Interface()
		return starlark.MakeInt(v.(int))
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := f.Interface()
		return starlark.MakeInt64(v.(int64))
	case reflect.Float32, reflect.Float64:
		v := f.Interface()
		return starlark.Float(v.(float64))
	}
	return starlark.None
}
