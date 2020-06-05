// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package schema

import (
	"fmt"
	"github.com/orivil/types"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

type Schema struct {
	Name        string       `json:"name,omitempty"`
	Model       string       `json:"model,omitempty"`
	Namespace   string       `json:"namespace,omitempty"`
	Ref         string       `json:"$ref,omitempty"`
	Type        JsonKind     `json:"type,omitempty"`
	Description string       `json:"description,omitempty"`
	Items       *Schema      `json:"items,omitempty"`
	Properties  Properties   `json:"properties,omitempty"`
	Validations *Validations `json:"validations,omitempty"`
}

type Models map[string]*Schema

func (ms Models) GetSchema(ref string) *Schema {
	return ms[ref]
}

func NewSchema(v interface{}) (*Schema, error) {
	rv := reflect.ValueOf(v)
	rt := indirectType(rv.Type())
	if rt.Kind() != reflect.Struct {
		return nil, fmt.Errorf("schema.NewSchemas: v must be struct or pointer of struct")
	} else {
		schema, err := valueToSchema(rv, nil, make(map[reflect.Type]struct{}))
		if err != nil {
			return nil, err
		}
		return schema, nil
	}
}

func (s *Schema) Models() Models {
	models := Models{}
	if s.Model != "" {
		models[s.Namespace+"."+s.Model] = s
	}
	for _, schema := range s.Properties {
		subs := schema.Models()
		for name, sub := range subs {
			models[name] = sub
		}
	}
	if s.Items != nil {
		subs := s.Items.Models()
		for name, sub := range subs {
			models[name] = sub
		}
	}
	return models
}

func (s *Schema) Valid(v interface{}) (info *Validations, err error) {
	return s.valid(reflect.ValueOf(v))
}

func (s *Schema) valid(v reflect.Value) (info *Validations, err error) {
	if s.Validations != nil && s.Type != Object {
		v = reflect.Indirect(v)
		//if t := v.Type(); s.GoType != t {
		//	return nil, fmt.Errorf("schema type: %s, got value type: %s", s.GoType, t)
		//}
		valid := v.IsValid() && !v.IsZero()
		if s.Validations.Required && !valid {
			return &Validations{Required: true}, nil
		}
		if valid {
			switch s.Type {
			case String, Number:
				var tv types.Value
				tv, err = types.GetValue(v.Interface())
				if err != nil {
					return nil, err
				}
				switch s.Type {
				case String:
					info, err = s.Validations.validString(tv.String())
					if err != nil || info != nil {
						return info, err
					}
				case Number:
					var num float64
					num, err = tv.Float64()
					if err != nil {
						return nil, err
					}
					info, err = s.Validations.validNumber(num)
					if err != nil || info != nil {
						return info, err
					}
				}
			case Array:
				info = s.Validations.validItemsLength(v.Len())
				if info != nil {
					return info, nil
				}
			}
		}
	}
	if s.Type == Object {
		v = indirectValue(v, true)
		vk := v.Kind()
		if vk == reflect.Struct {
			fs := getStructFields(v)
			fvs := make(map[string]reflect.Value, len(fs))
			for _, f := range fs {
				fvs[f.property] = f.fv
			}
			for _, schema := range s.Properties {
				fv := fvs[schema.Name]
				info, err = schema.valid(fv)
				if info != nil || err != nil {
					return info, err
				}
			}
		} else if vk == reflect.Map {
			for property, schema := range s.Properties {
				fv := v.MapIndex(reflect.ValueOf(property))
				info, err = schema.valid(fv)
				if info != nil || err != nil {
					return info, err
				}
			}
		}
	}
	return nil, nil
}

func (s *Schema) initValidation() {
	if s.Validations == nil {
		s.Validations = &Validations{}
	}
}

func (s *Schema) Requires(properties []string) *Schema {
	for _, property := range properties {
		for _, schema := range s.Properties {
			if schema.Name == property {
				schema.WithRequired(true)
				break
			}
		}
	}
	return s
}

func (s *Schema) Property(property string) *Schema {
	for _, schema := range s.Properties {
		if schema.Name == property {
			return schema
		}
	}
	return nil
}

func (s *Schema) WithTagOptions(st reflect.StructTag) error {
	if st != "" {
		if desc := st.Get(Description); desc != "" {
			s.WithDescription(desc)
		}
		optStr := st.Get(Tag)
		if optStr != "" {
			opts, err := parseTag(optStr)
			if err != nil {
				return &TagError{
					Tag: Tag,
					Err: err.Error(),
				}
			}
			if opts.Contains(OptionsRequired) {
				s.WithRequired(true)
			}
			if str := opts.GetValue(Enum); str != "" {
				elements := strings.Split(str, ",")
				var enum = make([]interface{}, len(elements))
				for i, elem := range elements {
					enum[i] = strings.TrimSpace(elem)
				}
				err = s.withEnum(enum)
				if err != nil {
					return &TagError{
						Tag: Tag + "." + Enum,
						Err: err.Error(),
					}
				}
			}
			var f64 float64
			if minNum := opts.GetValue(MinNum); minNum != "" {
				f64, err = strToFloat64(minNum)
				if err != nil {
					return &TagError{
						Tag: Tag + "." + MinNum,
						Err: err.Error(),
					}
				}
				s.WithMinNum(f64)
			}
			if maxNum := opts.GetValue(MaxNum); maxNum != "" {
				f64, err = strToFloat64(maxNum)
				if err != nil {
					return &TagError{
						Tag: Tag + "." + MaxNum,
						Err: err.Error(),
					}
				}
				s.WithMaxNum(f64)
			}
			if minExcNum := opts.GetValue(MinExcNum); minExcNum != "" {
				f64, err = strToFloat64(minExcNum)
				if err != nil {
					return &TagError{
						Tag: Tag + "." + MinExcNum,
						Err: err.Error(),
					}
				}
				s.WithMinExcNum(f64)
			}
			if maxExcNum := opts.GetValue(MaxExcNum); maxExcNum != "" {
				f64, err = strToFloat64(maxExcNum)
				if err != nil {
					return &TagError{
						Tag: Tag + "." + MaxExcNum,
						Err: err.Error(),
					}
				}
				s.WithMaxExcNum(f64)
			}
			var i int
			if minLen := opts.GetValue(MinLen); minLen != "" {
				i, err = types.String(minLen).Int()
				if err != nil {
					return &TagError{
						Tag: Tag + "." + MinLen,
						Err: err.Error(),
					}
				}
				s.WithMinLen(i)
			}
			if maxLen := opts.GetValue(MaxLen); maxLen != "" {
				i, err = types.String(maxLen).Int()
				if err != nil {
					return &TagError{
						Tag: Tag + "." + MaxLen,
						Err: err.Error(),
					}
				}
				s.WithMaxLen(i)
			}
			if minItems := opts.GetValue(MinItems); minItems != "" {
				i, err = types.String(minItems).Int()
				if err != nil {
					return &TagError{
						Tag: Tag + "." + MinItems,
						Err: err.Error(),
					}
				}
				s.WithMinItems(i)
			}
			if maxItems := opts.GetValue(MaxItems); maxItems != "" {
				i, err = types.String(maxItems).Int()
				if err != nil {
					return &TagError{
						Tag: Tag + "." + MaxItems,
						Err: err.Error(),
					}
				}
				s.WithMaxItems(i)
			}
			if pattern := opts.GetValue(Pattern); pattern != "" {
				err = s.withPattern(pattern)
				if err != nil {
					return &TagError{
						Tag: Tag + "." + Pattern,
						Err: err.Error(),
					}
				}
			}
		}
	}
	return nil
}

func (s *Schema) WithDescription(description string) *Schema {
	s.Description = description
	return s
}

func strToFloat64(s string) (float64, error) {
	return types.String(s).Float64()
}

func (s *Schema) WithRequired(required bool) *Schema {
	s.initValidation()
	s.Validations.Required = required
	return s
}

var patterns = &matchers{}

func (s *Schema) WithPattern(pattern string) *Schema {
	err := s.withPattern(pattern)
	if err != nil {
		panic(err)
	}
	return s
}

func (s *Schema) withPattern(pattern string) error {
	s.initValidation()
	_, err := patterns.Compile(pattern)
	if err != nil {
		return err
	}
	s.Validations.Pattern = pattern
	return nil
}

func (s *Schema) WithMaxNum(maxNum float64) *Schema {
	s.initValidation()
	s.Validations.MaxExcNum = nil
	s.Validations.MaxNum = &maxNum
	return s
}
func (s *Schema) WithMinNum(minNum float64) *Schema {
	s.initValidation()
	s.Validations.MinExcNum = nil
	s.Validations.MinNum = &minNum
	return s
}
func (s *Schema) WithMaxExcNum(maxExcNum float64) *Schema {
	s.initValidation()
	s.Validations.MaxNum = nil
	s.Validations.MaxExcNum = &maxExcNum
	return s
}
func (s *Schema) WithMinExcNum(minExcNum float64) *Schema {
	s.initValidation()
	s.Validations.MinNum = nil
	s.Validations.MinExcNum = &minExcNum
	return s
}
func (s *Schema) WithMaxLen(maxLen int) *Schema {
	s.initValidation()
	s.Validations.MaxLen = &maxLen
	return s
}
func (s *Schema) WithMinLen(minLen int) *Schema {
	s.initValidation()
	s.Validations.MinLen = &minLen
	return s
}
func (s *Schema) WithMaxItems(maxItems int) *Schema {
	s.initValidation()
	s.Validations.MaxItems = &maxItems
	return s
}
func (s *Schema) WithMinItems(minItems int) *Schema {
	s.initValidation()
	s.Validations.MinItems = &minItems
	return s
}
func (s *Schema) withEnum(enum []interface{}) error {
	s.initValidation()
	var elements []interface{}
	for _, e := range enum {
		elem, err := types.ToValue(reflect.String, e)
		if err != nil {
			return err
		}
		elements = append(elements, elem)
	}
	s.Validations.Enum = elements
	return nil
}
func (s *Schema) WithEnum(enum ...interface{}) *Schema {
	err := s.withEnum(enum)
	if err != nil {
		panic(err)
	}
	return s
}

type matchers struct {
	syncMap sync.Map
}

func (ms *matchers) Compile(pattern string) (*regexp.Regexp, error) {
	var reg *regexp.Regexp
	v, ok := ms.syncMap.Load(pattern)
	if !ok {
		var err error
		reg, err = regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		ms.syncMap.Store(pattern, reg)
	} else {
		reg = v.(*regexp.Regexp)
	}
	return reg, nil
}
