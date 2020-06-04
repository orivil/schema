# Golang Package For Describing JSON Object

* Note: This package is **NOT** implement the [json-schema](https://json-schema.org/)!

## Example
```go
package main

import (
	"fmt"
	."github.com/orivil/schema"
	"net/url"
)

func main() {
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
		"sex": []string{"3"},
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
	fmt.Println("enum:", info.Enum)

	// Output:
	// username: JayChou
	// password: ChouJay
	// sex: 3
	// enum: [1 2]
}
```