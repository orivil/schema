// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package schema_test

import (
	"fmt"
	. "github.com/orivil/schema"
	"net/url"
)

func ExampleUnmarshalUrl() {
	type params struct {
		Username string `json:"username" schema:"required; pattern:[\\w]{6,12}"`
		Password string `json:"password" schema:"required; pattern:[\\w]{6,12}"`
		Sex      *int   `json:"sex" schema:"required; enum:1,2" desc:"1-male, 2-female"`
	}
	schema, err := NewSchema(params{})
	if err != nil {
		panic(err)
	}
	urlValues := url.Values{
		"username": []string{"JayChou"},
		"password": []string{"ChouJay"},
		"sex":      []string{"3"},
	}
	ps := &params{}
	err = UnmarshalUrl(urlValues, ps)
	if err != nil {
		panic(err)
	}
	info, err := schema.Valid(ps)
	if err != nil {
		panic(err)
	}
	fmt.Println("username:", ps.Username)
	fmt.Println("password:", ps.Password)
	fmt.Println("sex:", *ps.Sex)
	fmt.Printf("field: %s, enum: %v", info.Field, info.Enum)

	// Output:
	// username: JayChou
	// password: ChouJay
	// sex: 3
	// field: sex, enum: [1 2]
}
