// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package schema

import (
	"fmt"
	"github.com/orivil/types"
	"reflect"
)

type Decoder interface {
	Schema() *Schema
}

var decoderType = reflect.TypeOf(new(Decoder)).Elem()

func valueToSchema(v reflect.Value, root *reflect.Type, existStructs map[reflect.Type]struct{}) (*Schema, error) {
	if !v.IsValid() {
		return nil, nil
	}
	if v.Type().Implements(decoderType) {
		return v.Interface().(Decoder).Schema(), nil
	}
	v = indirectValue(v, true)
	t := v.Type()
	schema := &Schema{Type: GoToJSONType(t)}
	k := t.Kind()
	switch k {
	case reflect.Interface:
		return valueToSchema(v.Elem(), root, existStructs)
	case reflect.Slice, reflect.Array:
		if schema.Type != File {
			ln := v.Len()
			var (
				items *Schema
				err   error
			)
			if ln == 0 { // nil slice
				items, err = valueToSchema(reflect.New(t.Elem()), root, existStructs)
			} else {
				for i := 0; i < ln-1; i++ {
					pre := indirectType(v.Index(i).Type())
					next := indirectType(v.Index(i + 1).Type())
					if pre != next {
						return nil, fmt.Errorf("slice or array element type must be unique, got %s, and %s", pre, next)
					}
				}
				items, err = valueToSchema(v.Index(0), root, existStructs)
			}
			if err != nil {
				return nil, err
			}
			if items != nil {
				schema.Items = items
			}
		}
	case reflect.Struct:
		if _, ok := existStructs[t]; ok {
			return &Schema{Ref: t.PkgPath() + "." + t.Name()}, nil
		} else {
			schema.Model = t.Name()
			schema.Namespace = t.PkgPath()
			existStructs[t] = struct{}{}
		}
		schema.Properties = Properties{}
		fields := getStructFields(v)
		for _, field := range fields {
			if ignore := isFieldIgnored(field.ft.Tag); ignore {
				continue
			}
			fs, err := valueToSchema(field.fv, root, existStructs)
			if err != nil {
				return nil, err
			}
			if fs != nil {
				fs.Field = field.ft.Name
				err = fs.WithTagOptions(field.ft.Tag)
				if err != nil {
					return nil, err
				}
				schema.Properties[field.property] = fs
			}
		}
	case reflect.Map:
		keys := v.MapKeys()
		for _, key := range keys {
			mv := v.MapIndex(key)
			ms, err := valueToSchema(mv, root, existStructs)
			if err != nil {
				return nil, err
			}
			if ms != nil {
				if key.CanInterface() {
					// get key string type
					var pv types.Value
					pv, err = types.GetValue(key.Interface())
					if err != nil {
						return nil, err
					} else {
						property := pv.String()
						ms.Field = property
						schema.Properties[property] = ms
					}
				}
			}
		}
	}
	return schema, nil
}

func indirectValue(v reflect.Value, newEmpty bool) reflect.Value {
	if v.Kind() == reflect.Ptr {
		if newEmpty && v.IsNil() {
			v = reflect.New(v.Type().Elem())
		}
		return indirectValue(v.Elem(), newEmpty)
	}
	return v
}

func indirectType(t reflect.Type) reflect.Type {
	switch t.Kind() {
	case reflect.Ptr:
		return indirectType(t.Elem())
	default:
		return t
	}
}

type structField struct {
	property string
	fv       reflect.Value
	ft       reflect.StructField
}

// get struct fields and merge anonymous fields
func getStructFields(v reflect.Value) []*structField {
	v = indirectValue(v, true)
	if v.Kind() == reflect.Struct {
		fieldNum := v.NumField()
		t := v.Type()
		var fields []*structField
		exists := make(map[string]struct{})
		var subFields []*structField
		for i := 0; i < fieldNum; i++ {
			fv := v.Field(i)
			ft := t.Field(i)
			if ft.Anonymous {
				subFields = append(subFields, getStructFields(fv)...)
			} else {
				property := getFieldName(ft.Tag)
				if property == "" {
					property = ft.Name
				}
				sf := &structField{
					property: property,
					fv:       indirectValue(fv, true),
					ft:       ft,
				}
				fields = append(fields, sf)
				exists[ft.Name] = struct{}{}
			}
		}
		for _, field := range subFields {
			if _, ok := exists[field.ft.Name]; !ok {
				fields = append(fields, field)
				exists[field.ft.Name] = struct{}{}
			}
		}
		return fields
	}
	return nil
}
