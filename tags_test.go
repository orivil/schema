// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"testing"
)

func TestTagParsing(t *testing.T) {
	type result struct {
		key, value string
	}
	type cases struct {
		tag      string
		results  []result
		contains []string
	}
	var testCases = cases{
		tag: "required;min:18;max:20;unique",
		results: []result{
			{"required", ""},
			{"min", "18"},
			{"max", "20"},
			{"unique", ""},
		},
		contains: []string{"required", "unique", "min", "max"},
	}
	opts, err := parseTag(testCases.tag)
	if err != nil {
		t.Fatal(err)
	}
	for _, res := range testCases.results {
		value := opts.GetValue(res.key)
		if value != res.value {
			t.Errorf("need: %s\ngot: %s\n", res.value, value)
		}
	}
	for _, contain := range testCases.contains {
		if !opts.Contains(contain) {
			t.Errorf("not contain %s\n", contain)
		}
	}
}
