// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package confi

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/naoina/toml"
	"github.com/naoina/toml/ast"
)

// Read TOML into the configuration.
func Read(r io.Reader, config interface{}) error {
	return read(r, config, false)
}

func read(r io.Reader, config interface{}, ignoreUnknown bool) (err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	table, err := toml.Parse(data)
	if err != nil {
		return
	}

	defer func() {
		err = asError(recover())
	}()

	setFields(config, "", table.Fields, ignoreUnknown)
	return
}

// ReadFile containing TOML into the configuration.
func ReadFile(filename string, config interface{}) error {
	return readFile(filename, config, false)
}

func readFile(filename string, config interface{}, ignoreUnknown bool) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	return read(f, config, ignoreUnknown)
}

// ReadFileIfExists is a lenient alternative to the ReadFile method..  No error
// is returned if the file doesn't exist.
func ReadFileIfExists(filename string, config interface{}) error {
	return readFileIfExists(filename, config, false)
}

func readFileIfExists(filename string, config interface{}, ignoreUnknown bool) (err error) {
	err = readFile(filename, config, ignoreUnknown)
	if err != nil && os.IsNotExist(err) {
		err = nil
	}
	return
}

func setFields(config interface{}, path string, fields map[string]interface{}, ignoreUnknown bool) {
	for k, v := range fields {
		p := k
		if path != "" {
			p = path + "." + k
		}

		switch x := v.(type) {
		case *ast.KeyValue:
			var s string

			switch y := x.Value.(type) {
			case *ast.Array:
				s = x.Value.Source()
			case *ast.Boolean:
				s = y.Value
			case *ast.Float:
				s = y.Value
			case *ast.Integer:
				s = y.Value
			case *ast.String:
				s = y.Value
			default:
				panic(fmt.Errorf("%s: type not supported: %#v", p, x.Value))
			}

			if ignoreUnknown {
				func() {
					defer func() {
						if x := recover(); x != nil {
							if _, ok := x.(unknownKeyError); !ok {
								panic(x)
							}
						}
					}()
					MustSetFromString(config, p, s)
				}()
			} else {
				MustSetFromString(config, p, s)
			}

		case *ast.Table:
			setFields(config, p, x.Fields, ignoreUnknown)

		default:
			panic(fmt.Errorf("%s: unknown value type: %#v", p, v))
		}
	}
}

// Write the configuration as TOML.
func Write(w io.Writer, config interface{}) error {
	return toml.NewEncoder(w).Encode(sanitizeContainer(make(map[string]interface{}), reflect.ValueOf(config).Elem()))
}

// WriteFile containing the configuration as TOML.
func WriteFile(filename string, config interface{}) (err error) {
	b := bytes.NewBuffer(nil)

	err = toml.NewEncoder(b).Encode(sanitizeContainer(make(map[string]interface{}), reflect.ValueOf(config).Elem()))
	if err != nil {
		return
	}

	return ioutil.WriteFile(filename, b.Bytes(), 0666)
}

func sanitizeContainer(sane map[string]interface{}, node reflect.Value) map[string]interface{} {
	switch node.Kind() {
	case reflect.Map:
		return sanitizeMap(sane, node)

	case reflect.Struct:
		return sanitizeStruct(sane, node)

	default:
		panic("must be a struct or a map")
	}
}

func sanitizeMap(sane map[string]interface{}, node reflect.Value) map[string]interface{} {
	for _, key := range reflectMapKeyStrings(node) {
		value := node.MapIndex(reflect.ValueOf(key))

		if value.Kind() == reflect.Interface {
			value = value.Elem()
		}

		if x := sanitizeValue(sane, value, false); x != nil {
			sane[key] = x
		}
	}

	return sane
}

func sanitizeStruct(sane map[string]interface{}, node reflect.Value) map[string]interface{} {
	for i := 0; i < node.Type().NumField(); i++ {
		value := node.Field(i)
		if !value.CanInterface() {
			continue
		}

		field := node.Type().Field(i)

		if x := sanitizeValue(sane, value, field.Anonymous); x != nil {
			sane[strings.ToLower(field.Name)] = x
		}
	}

	return sane
}

func sanitizeValue(sane map[string]interface{}, value reflect.Value, anonymous bool) (x interface{}) {
	switch value.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
		x = value.Interface()

	case reflect.Int64:
		x = value.Interface()
		if d, ok := x.(time.Duration); ok {
			x = d.String()
		}

	case reflect.Slice:
		switch value.Type().Elem().Kind() {
		case reflect.String:
			x = value.Interface()
		}

	case reflect.Map:
		if s := sanitizeMap(make(map[string]interface{}), value); len(s) > 0 {
			x = s
		}

	case reflect.Ptr:
		if value.IsNil() {
			break
		}
		if value.Type().Elem().Kind() != reflect.Struct {
			break
		}
		value = value.Elem()
		fallthrough

	case reflect.Struct:
		if anonymous {
			sane = sanitizeStruct(sane, value)
		} else if s := sanitizeStruct(make(map[string]interface{}), value); len(s) > 0 {
			x = s
		}
	}

	return
}
