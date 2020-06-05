// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package schema

import (
	"github.com/orivil/types"
	"regexp"
)

type Validations struct {
	Required  bool     `json:"required,omitempty"`
	Pattern   string   `json:"pattern,omitempty"`
	MaxItems  *int     `json:"maxItems,omitempty"`
	MinItems  *int     `json:"minItems,omitempty"`
	MaxLen    *int     `json:"maxLen,omitempty"`
	MinLen    *int     `json:"minLen,omitempty"`
	MaxNum    *float64 `json:"maxNum,omitempty"`
	MinNum    *float64 `json:"minNum,omitempty"`
	MaxExcNum *float64 `json:"maxExcNum,omitempty"`
	MinExcNum *float64 `json:"minExcNum,omitempty"`
	Enum      []string `json:"enum,omitempty"`
}

func (vs *Validations) validItemsLength(length int) *Validations {
	if (vs.MinLen != nil && *vs.MinLen > length) || (vs.MaxLen != nil && *vs.MaxLen < length) {
		return &Validations{MinItems: vs.MinItems, MaxItems: vs.MaxItems}
	}
	return nil
}

func (vs *Validations) validNumber(num float64) (info *Validations, err error) {
	if vs.Enum != nil {
		exist := false
		var f64 float64
		for _, enum := range vs.Enum {
			str := types.String(enum)
			f64, err = str.Float64()
			if err != nil {
				return nil, err
			}
			if f64 == num {
				exist = true
				break
			}
		}
		if !exist {
			return &Validations{Enum: vs.Enum}, nil
		}
	}
	if vs.MinNum != nil {
		if *vs.MinNum > num {
			return &Validations{MinNum: vs.MinNum}, nil
		}
	}
	if vs.MaxNum != nil {
		if *vs.MaxNum < num {
			return &Validations{MaxNum: vs.MaxNum}, nil
		}
	}
	if vs.MinExcNum != nil {
		if *vs.MinExcNum >= num {
			return &Validations{MinExcNum: vs.MinExcNum}, nil
		}
	}
	if vs.MaxExcNum != nil {
		if *vs.MaxExcNum < num {
			return &Validations{MaxExcNum: vs.MaxExcNum}, nil
		}
	}
	return nil, nil
}

func (vs *Validations) validString(str string) (info *Validations, err error) {
	if vs.Enum != nil {
		exist := false
		for _, e := range vs.Enum {
			if e == str {
				exist = true
				break
			}
		}
		if !exist {
			return &Validations{Enum: vs.Enum}, nil
		}
	}
	if vs.Pattern != "" {
		var matcher *regexp.Regexp
		matcher, err = patterns.Compile(vs.Pattern)
		if err != nil {
			return nil, err
		}
		if !matcher.MatchString(str) {
			return &Validations{Pattern: vs.Pattern}, nil
		}
	}
	length := len(str)
	if vs.MinLen != nil && *vs.MinLen > length {
		return &Validations{MinLen: vs.MinLen}, nil
	}
	if vs.MaxLen != nil && *vs.MaxLen < length {
		return &Validations{MaxLen: vs.MaxLen}, nil
	}
	return nil, nil
}
