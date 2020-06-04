// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package schema

import (
	"errors"
	"fmt"
	"strings"
)

const (
	Tag         = "schema"
	Description = "desc"
	Enum        = "enum"
	MaxNum      = "maxNum"
	MinNum      = "minNum"
	MinExcNum   = "minExcNum"
	MaxExcNum   = "maxExcNum"
	MinLen      = "minLen"
	MaxLen      = "maxLen"
	MinItems    = "minItems"
	MaxItems    = "maxItems"
	Pattern     = "pattern"
)

const (
	OptionsRequired = "required"
)

type TagError struct {
	Tag string
	Err string
}

func (t *TagError) Error() string {
	return fmt.Sprintf("schema tag [%s] got error: %s", t.Tag, t.Err)
}

type tagOptions map[string]string

func parseTag(tag string) (tagOptions, error) {
	options := strings.Split(tag, ";")
	var ts = make(tagOptions)
	for _, option := range options {
		option = strings.TrimSpace(option)
		if len(option) > 0 {
			ln := len(option)
			sepIdx := strings.Index(option, ":")
			if sepIdx != -1 {
				if sepIdx == 0 {
					return nil, errors.New("invalid options key")
				} else {
					key := option[0:sepIdx]
					ts[key] = ""
					if sepIdx < ln-1 {
						ts[key] = option[sepIdx+1:]
					}
				}
			} else {
				ts[option] = ""
			}
		}
	}
	return ts, nil
}

func (opts tagOptions) Contains(optionName string) bool {
	_, ok := opts[optionName]
	return ok
}

func (opts tagOptions) GetValue(optionName string) string {
	return opts[optionName]
}
