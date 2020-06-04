// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package schema_test

import (
	"github.com/orivil/schema"
	"net/url"
	"testing"
)

func TestUnmarshalUrlValues(t *testing.T) {
	type Anonymous struct {
		Bool *bool `json:"bool"`
	}
	type Params struct {
		Str        string   `json:"str"`
		Int        *int     `json:"int"`
		Int8       int8     `json:"int_8"`
		SliceInt   []int    `json:"slice_int"`
		SliceStr   []string `json:"slice_str"`
		*Anonymous `json:"-"`
	}
	values := url.Values{
		"str":       []string{"Nina"},
		"int":       []string{"64"},
		"int_8":     []string{"8"},
		"slice_int": []string{"1", "2"},
		"slice_str": []string{"1", "2"},
		"bool":      []string{"true"},
	}
	params := &Params{}
	err := schema.UnmarshalUrl(values, params)
	if err != nil {
		t.Fatal(err)
	}
	got := jsonStr(params)
	need := jsonStr(Params{
		Str:       "Nina",
		Int:       newInt(64),
		Int8:      8,
		SliceInt:  []int{1, 2},
		SliceStr:  []string{"1", "2"},
		Anonymous: &Anonymous{Bool: newBool(true)},
	})
	if got != need {
		t.Fatalf("need: %s\ngot: %s", need, got)
	}
}

// 973 ns/op
func BenchmarkUnmarshalUrlValues(b *testing.B) {
	type Params struct {
		Str string `json:"str"`
		Int *int   `json:"int"`
	}
	//s, err := schema.NewSchema(&Params{})
	//if err != nil {
	//	b.Fatal(err)
	//}
	values := url.Values{
		"str": []string{"Nina"},
		"int": []string{"64"},
	}
	for i := 0; i < b.N; i++ {
		params := &Params{}
		err := schema.UnmarshalUrl(values, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func newInt(i int) *int    { return &i }
func newBool(b bool) *bool { return &b }
