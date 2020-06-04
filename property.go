// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package schema

import (
	"reflect"
	"strings"
)

type Properties map[string]*Schema

func getFieldName(tag reflect.StructTag) string {
	v := tag.Get("json")
	if idx := strings.Index(v, ","); idx != -1 {
		return v[:idx]
	} else {
		return v
	}
}

func isFieldIgnored(tag reflect.StructTag) bool {
	return tag.Get("json") == "-"
}
