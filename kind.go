// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package schema

import (
	"io/ioutil"
	"mime/multipart"
	"reflect"
)

const (
	Invalid JsonKind = "invalid"
	Bool    JsonKind = "Boolean"
	Number  JsonKind = "Number"
	Array   JsonKind = "Array"
	Object  JsonKind = "Object"
	File    JsonKind = "File"
	String  JsonKind = "String"
)

type JsonKind string

type FileInterface interface {
	Read(header multipart.FileHeader) error
}

type FileData []byte

func (f *FileData) Read(header multipart.FileHeader) (err error) {
	var fs multipart.File
	defer func() {
		if fs != nil {
			fs.Close()
		}
	}()
	fs, err = header.Open()
	if err != nil {
		return err
	}
	var data []byte
	data, err = ioutil.ReadAll(fs)
	if err != nil {
		return err
	}
	*f = data
	return nil
}

var fileDataType = reflect.TypeOf(new(FileInterface)).Elem()

var reflectKinds = map[reflect.Kind]JsonKind{
	reflect.Invalid:       Invalid,
	reflect.Bool:          Bool,
	reflect.Int:           Number,
	reflect.Int8:          Number,
	reflect.Int16:         Number,
	reflect.Int32:         Number,
	reflect.Int64:         Number,
	reflect.Uint:          Number,
	reflect.Uint8:         Number,
	reflect.Uint16:        Number,
	reflect.Uint32:        Number,
	reflect.Uint64:        Number,
	reflect.Uintptr:       Invalid,
	reflect.Float32:       Number,
	reflect.Float64:       Number,
	reflect.Complex64:     Invalid,
	reflect.Complex128:    Invalid,
	reflect.Array:         Array,
	reflect.Chan:          Invalid,
	reflect.Func:          Invalid,
	reflect.Interface:     Invalid,
	reflect.Map:           Object,
	reflect.Ptr:           Invalid,
	reflect.Slice:         Array,
	reflect.String:        String,
	reflect.Struct:        Object,
	reflect.UnsafePointer: Invalid,
}

func GoToJSONType(t reflect.Type) JsonKind {
	t = indirectType(t)
	if reflect.New(t).Type().ConvertibleTo(fileDataType) {
		return File
	} else {
		return reflectKinds[t.Kind()]
	}
}
