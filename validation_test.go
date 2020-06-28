// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package schema_test

import (
	. "github.com/orivil/schema"
	"reflect"
	"testing"
)

func TestValidate(t *testing.T) {
	type Anonymous struct {
		SStr  []string `json:"s_str"`
		SInt  []int8   `json:"s_int"`
		SBool []bool   `json:"s_bool"`
	}
	type model struct {
		Str  string `json:"str"`
		Int  int    `json:"int"`
		Int8 int8   `json:"int8"`
		*Anonymous
	}
	type testCase struct {
		field string
		tag   string
		v     model
		info  *Validations
	}
	var testCases = []testCase{
		{"str", `schema:"required"`, model{}, &Validations{Field: "str", Required: true}},
		{"str", `schema:"minLen:2"`, model{}, nil},
		{"str", `schema:"minLen:2"`, model{Str: "1"}, &Validations{Field: "str", MinLen: newInt(2)}},
		{"str", `schema:"maxLen:2"`, model{}, nil},
		{"str", `schema:"maxLen:2"`, model{Str: "123"}, &Validations{Field: "str", MaxLen: newInt(2)}},
		{"str", `schema:"pattern:^Jay"`, model{}, nil},
		{"str", `schema:"pattern:^Jay"`, model{Str: "Chou Jay"}, &Validations{Field: "str", Pattern: "^Jay"}},
		{"int", `schema:"required"`, model{}, &Validations{Field: "int", Required: true}},
		{"int", `schema:"minNum:2"`, model{}, nil},
		{"int", `schema:"minNum:2"`, model{Int: 1}, &Validations{Field: "int", MinNum: newFloat(2)}},
		{"int", `schema:"maxNum:2"`, model{Int: 3}, &Validations{Field: "int", MaxNum: newFloat(2)}},
		{"int8", `schema:"required"`, model{}, &Validations{Field: "int8", Required: true}},
		{"s_str", `schema:"required"`, model{}, &Validations{Field: "s_str", Required: true}},
		{"s_str", `schema:"minItems:2"`, model{Anonymous: &Anonymous{SStr: []string{"1", "2"}}}, nil},
		{"s_str", `schema:"maxItems:2"`, model{Anonymous: &Anonymous{SStr: []string{"1", "2"}}}, nil},
		{"s_str", `schema:"minItems:2;maxItems:3"`, model{Anonymous: &Anonymous{SStr: []string{"1", "2", "3", "4"}}}, &Validations{Field: "s_str", MinItems: newInt(2), MaxItems: newInt(3)}},
		{"s_int", `schema:"required"`, model{}, &Validations{Field: "s_int", Required: true}},
	}
	for _, tc := range testCases {
		schema, err := NewSchema(&model{})
		if err != nil {
			t.Fatal(err)
		}
		err = schema.Property(tc.field).WithTagOptions(reflect.StructTag(tc.tag))
		if err != nil {
			t.Fatal(err)
		}
		var info *Validations
		info, err = schema.Valid(tc.v)
		if err != nil {
			t.Fatal(err)
		}
		if info != nil {
			got := jsonStr(info)
			need := jsonStr(tc.info)
			if got != need {
				t.Errorf("need: %s, got: %s\n", need, got)
			}
		}
	}
}

func newFloat(f float64) *float64 {
	return &f
}
