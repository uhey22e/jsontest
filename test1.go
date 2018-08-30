package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/volatiletech/null"
)

type NullStr null.String

func (s NullStr) MarshalJSON() ([]byte, error) {
	if s.Valid {
		b, e := json.Marshal(s.String)
		return b, e
	} else {
		return []byte{}, nil
	}
}

func newNullStr(s null.String) NullStr {
	return NullStr{
		String: s.String,
		Valid:  s.Valid,
	}
}

type testStruct1 struct {
	TypeInt     int     `json:"typeInt"`
	TypeIntPtr  *int    `json:"typeIntPtr"`
	TypeString  string  `json:"typeString"`
	TypeStringP *string `json:"typeStringP"`
}

type testStruct2 struct {
	TypeInt     int     `json:"typeInt"`
	TypeIntPtr  *int    `json:"typeIntPtr,omitempty"`
	TypeString  string  `json:"typeString"`
	TypeStringP *string `json:"typeStringP,omitempty"`
}

type testStruct3 struct {
	TypeInt     int     `json:"typeInt"`
	TypeIntPtr  *int    `json:"typeIntPtr,string"`
	TypeString  string  `json:"typeString"`
	TypeStringP *string `json:"typeStringP,string"`
}

type testStruct4 struct {
	TypeInt     int         `json:"typeInt"`
	TypeNInt    null.Int    `json:"typeIntPtr"`
	TypeString  string      `json:"typeString"`
	TypeNString null.String `json:"typeStringP"`
}

type TestStruct5 struct {
	TypeInt     int      `json:"typeInt"`
	TypeNInt    null.Int `json:"typeIntPtr"`
	TypeString  string   `json:"typeString"`
	TypeNString NullStr  `json:"typeStringP"`
}

func test1() {
	testjson := `
	{
		"typeInt": 123,
		"typeIntPtr": null,
		"typeString": "hogehoge",
		"typeStringP": null,
	}`

	// decode
	d1 := &testStruct1{}
	json.NewDecoder(strings.NewReader(testjson)).Decode(d1)
	fmt.Println(fmt.Sprintf("%#v", d1))

	// encode
	d2 := &testStruct1{
		TypeInt:     123,
		TypeIntPtr:  nil,
		TypeString:  "hogehoge",
		TypeStringP: nil,
	}
	buf2 := new(bytes.Buffer)
	json.NewEncoder(buf2).Encode(d2)
	fmt.Println(fmt.Sprintf("%s", buf2.String()))

	// encode with 'omitempty' tag
	d3 := &testStruct2{
		TypeInt:     123,
		TypeIntPtr:  nil,
		TypeString:  "hogehoge",
		TypeStringP: nil,
	}
	buf3 := new(bytes.Buffer)
	json.NewEncoder(buf3).Encode(d3)
	fmt.Println(fmt.Sprintf("%s", buf3.String()))

	// encode with 'string' tag
	d4 := &testStruct3{
		TypeInt:     123,
		TypeIntPtr:  nil,
		TypeString:  "hogehoge",
		TypeStringP: nil,
	}
	buf4 := new(bytes.Buffer)
	json.NewEncoder(buf4).Encode(d4)
	fmt.Println(fmt.Sprintf("%s", buf4.String()))

	// encode with null package
	d5 := &testStruct4{
		TypeInt:     123,
		TypeNInt:    null.IntFromPtr(nil),
		TypeString:  "fugafuga",
		TypeNString: null.StringFromPtr(nil),
	}
	buf5 := new(bytes.Buffer)
	json.NewEncoder(buf5).Encode(d5)
	fmt.Println(fmt.Sprintf("%s", buf5.String()))

	// encode with null package
	d6 := &TestStruct5{
		TypeInt:     123,
		TypeNInt:    null.IntFromPtr(nil),
		TypeString:  "encode with typed struct",
		TypeNString: newNullStr(null.StringFromPtr(nil)),
	}
	buf6 := new(bytes.Buffer)
	json.NewEncoder(buf6).Encode(d6)
	fmt.Println(fmt.Sprintf("%s", buf6.String()))

}
