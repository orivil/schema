// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package schema

import (
	"fmt"
	"github.com/orivil/types"
	"net/url"
	"reflect"
)

type UrlUnmarshaler interface {
	UnmarshalUrl(vs url.Values) error
}

var urlUnmarshalerType = reflect.TypeOf(new(UrlUnmarshaler)).Elem()

func UnmarshalUrl(values url.Values, v interface{}) error {
	rv := reflect.ValueOf(v)
	return unmarshalUrl(values, &rv)
}

func unmarshalUrl(values url.Values, rv *reflect.Value) error {
	if rv.Type().Implements(urlUnmarshalerType) {
		return rv.Interface().(UrlUnmarshaler).UnmarshalUrl(values)
	}
	irv := reflect.Indirect(*rv)
	ik := irv.Kind()
	if ik == reflect.Struct {
		numFields := irv.NumField()
		irt := irv.Type()
		for i := 0; i < numFields; i++ {
			fv := irv.Field(i)
			if fv.CanSet() {
				ft := irt.Field(i)
				if !isFieldIgnored(ft.Tag) {
					if ft.Anonymous {
						if ft.Type.Kind() == reflect.Ptr {
							var setV reflect.Value
							if fv.IsNil() {
								nv := reflect.New(ft.Type.Elem())
								fv.Set(nv)
								setV = nv.Elem()
							} else {
								setV = fv.Elem()
							}
							err := unmarshalUrl(values, &setV)
							if err != nil {
								return err
							}
						}
					} else {
						property := getFieldName(ft.Tag)
						if property == "" {
							property = ft.Name
						}
						if vs := values[property]; len(vs) > 0 {
							err := setUrlValue(vs, &fv)
							if err != nil {
								return err
							}
						}
					}
				}
			}
		}
	} else {
		return fmt.Errorf("only support struct or pointer of struct, got %s", ik)
	}
	return nil
}

func setUrlValue(values []string, vp *reflect.Value) error {
	var v = *vp
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			nv := reflect.New(v.Type().Elem())
			v.Set(nv)
			v = nv.Elem()
		} else {
			v = v.Elem()
		}
	}
	it := v.Type()
	ik := it.Kind()
	if ik == reflect.Slice {
		vs := make([]interface{}, len(values))
		for i, value := range values {
			vs[i] = value
		}
		i, err := types.ToSlice(it.Elem().Kind(), vs)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(i))
	} else if ik == reflect.Struct {
		urlValues, err := url.ParseQuery(values[0])
		if err != nil {
			return err
		}
		return unmarshalUrl(urlValues, vp)
	} else {
		i, err := types.ToValue(ik, values[0])
		if err != nil {
			return err
		}
		var nv reflect.Value
		if v.Kind() == reflect.Ptr {
			nv = reflect.New(it)
			nv.Elem().Set(reflect.ValueOf(i))
		} else {
			nv = reflect.ValueOf(i)
		}
		v.Set(nv)
	}
	return nil
}
