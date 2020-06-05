// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package schema_test

import (
	"encoding/json"
	"github.com/orivil/schema"
	"testing"
)

type A struct {
	F1  string   `json:"f01" schema:"required;pattern:^J" desc:"user name"`
	F2  string   `json:"f02" schema:"minLen:10;maxLen:12"`
	F3  *int     `json:"f03"`
	F4  int32    `json:"f04" schema:"enum:1,2,3"`
	F5  int64    `json:"f05" schema:"minNum:16;maxExcNum:18"`
	F6  float32  `json:"f06" schema:"minExcNum:16;maxNum:18"`
	F7  float64  `json:"f07"`
	F8  bool     `json:"f08"`
	F9  []int    `json:"f09"`
	F10 []string `json:"f10" schema:"minLen:10;maxLen:12"`
	*B
}

type B struct {
	F1  *int32  `json:"f01"` // covered by A
	F11 *string `json:"f11" schema:"required"`
	F12 *C      `json:"f12"`
}

type C struct {
	F13 bool              `json:"f13" schema:"required"`
	F14 schema.FileData   `json:"f14"`
	F15 []schema.FileData `json:"f15"`
	F16 *A                `json:"f16"`
	F17 *B                `json:"f17"`
}

func TestNewSchemas(t *testing.T) {
	s, err := schema.NewSchema(&A{})
	if err != nil {
		panic(err)
	}
	got := jsonStr(s)
	need :=
		`{
	"model": "A",
	"namespace": "github.com/orivil/schema_test",
	"type": "Object",
	"properties": [
		{
			"name": "f01",
			"type": "String",
			"description": "user name",
			"validations": {
				"required": true,
				"pattern": "^J"
			}
		},
		{
			"name": "f02",
			"type": "String",
			"validations": {
				"maxLen": 12,
				"minLen": 10
			}
		},
		{
			"name": "f03",
			"type": "Number"
		},
		{
			"name": "f04",
			"type": "Number",
			"validations": {
				"enum": [
					"1",
					"2",
					"3"
				]
			}
		},
		{
			"name": "f05",
			"type": "Number",
			"validations": {
				"minNum": 16,
				"maxExcNum": 18
			}
		},
		{
			"name": "f06",
			"type": "Number",
			"validations": {
				"maxNum": 18,
				"minExcNum": 16
			}
		},
		{
			"name": "f07",
			"type": "Number"
		},
		{
			"name": "f08",
			"type": "Boolean"
		},
		{
			"name": "f09",
			"type": "Array",
			"items": {
				"type": "Number"
			}
		},
		{
			"name": "f10",
			"type": "Array",
			"items": {
				"type": "String"
			},
			"validations": {
				"maxLen": 12,
				"minLen": 10
			}
		},
		{
			"name": "f11",
			"type": "String",
			"validations": {
				"required": true
			}
		},
		{
			"name": "f12",
			"model": "C",
			"namespace": "github.com/orivil/schema_test",
			"type": "Object",
			"properties": [
				{
					"name": "f13",
					"type": "Boolean",
					"validations": {
						"required": true
					}
				},
				{
					"name": "f14",
					"type": "File"
				},
				{
					"name": "f15",
					"type": "Array",
					"items": {
						"type": "File"
					}
				},
				{
					"name": "f16",
					"$ref": "github.com/orivil/schema_test.A"
				},
				{
					"name": "f17",
					"model": "B",
					"namespace": "github.com/orivil/schema_test",
					"type": "Object",
					"properties": [
						{
							"name": "f01",
							"type": "Number"
						},
						{
							"name": "f11",
							"type": "String",
							"validations": {
								"required": true
							}
						},
						{
							"name": "f12",
							"$ref": "github.com/orivil/schema_test.C"
						}
					]
				}
			]
		}
	]
}`
	if got != need {
		t.Fatalf("need: %s, got: %s", need, got)
	}
}

func jsonStr(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(data)
}
