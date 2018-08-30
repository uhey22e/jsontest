package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/volatiletech/null"
)

type NullStringMod struct {
	null.String
}

type SampleStruct struct {
	RequiredInt     int
	NullableInt1    null.Int
	NullableInt2    null.Int
	RequiredString  string
	NullableString1 null.String
	NullableString2 null.String
	NullableString3 NullStringMod
}

type DBStruct struct {
	RequiredInt     int
	NullableInt1    null.Int
	NullableInt2    null.Int
	RequiredString  string
	NullableString1 null.String
	NullableString2 null.String
	NullableString3 null.String
	OtherColumn     null.Time
}

func (o SampleStruct) MarshalJSON() ([]byte, error) {
	type Alias SampleStruct
	if !(o.NullableString1.Valid) {
		o.NullableString1.SetValid("")
	}
	if !(o.NullableString2.Valid) {
		o.NullableString2.SetValid("")
	}
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(o),
	})
}

func (o *SampleStruct) UnmarshalJSON(b []byte) error {
	type Alias SampleStruct
	d := Alias{}
	json.Unmarshal(b, &d)
	o.RequiredInt = d.RequiredInt
	o.NullableInt1 = d.NullableInt1
	o.NullableInt2 = d.NullableInt2
	o.RequiredString = d.RequiredString
	if d.NullableString1.String == "" {
		o.NullableString1 = null.NewString("", false)
	} else {
		o.NullableString1 = d.NullableString1
	}
	o.NullableString2 = d.NullableString2
	o.NullableString3 = d.NullableString3
	return nil
}

func (o NullStringMod) MarshalJSON() ([]byte, error) {
	type Alias NullStringMod
	if !o.Valid {
		o.SetValid("")
	}
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(o),
	})
}

func (o *NullStringMod) UnmarshalJSON(b []byte) error {
	var d string
	if err := json.Unmarshal(b, &d); err != nil {
		return err
	}
	if d == "" {
		o.Valid = false
	} else {
		o.SetValid(d)
	}
	return nil
}

func test2() {
	// encode test
	var buf bytes.Buffer
	data := &SampleStruct{
		RequiredInt:     123,
		NullableInt1:    null.NewInt(100, false),
		NullableInt2:    null.NewInt(100, true),
		RequiredString:  "required",
		NullableString1: null.String{},
		NullableString2: null.StringFrom("nullable"),
	}
	json.NewEncoder(&buf).Encode(data)
	fmt.Println(buf.String())

	// decode test
	req := `
	{
		"requiredInt": 1,
		"nullableInt1": 123,
		"nullableInt2": null,
		"requiredString": "required",
		"nullableString1": "nullable",
		"nullableString2": "",
		"nullableString3": ""
	}`
	tmp := &SampleStruct{}
	json.NewDecoder(strings.NewReader(req)).Decode(tmp)
	fmt.Println(fmt.Sprintf("%#v", tmp))

	fmt.Println(fmt.Sprintf("%#v", toDBStruct(tmp)))
}

func toDBStruct(d *SampleStruct) *DBStruct {
	s := &DBStruct{}
	s.RequiredInt = d.RequiredInt
	s.NullableInt1 = d.NullableInt1
	s.NullableInt2 = d.NullableInt2
	s.RequiredString = d.RequiredString
	s.NullableString1 = d.NullableString1
	// 詰替時にvalid check
	if d.NullableString2.Valid && d.NullableString2.String == "" {
		s.NullableString2 = null.NewString("", false)
	} else {
		s.NullableString2 = d.NullableString2
	}
	// embedded struct
	s.NullableString3 = null.StringFromPtr(d.NullableString3.Ptr())
	return s
}

func fromDBStruct(d *DBStruct) *SampleStruct {
	s := &SampleStruct{}
	s.RequiredInt = d.RequiredInt
	s.NullableInt1 = d.NullableInt1
	s.NullableInt2 = d.NullableInt2
	s.RequiredString = d.RequiredString
	s.NullableString1 = d.NullableString1
	// 詰替時にvalid check
	if !d.NullableString2.Valid {
		s.NullableString2 = null.StringFrom("")
	} else {
		s.NullableString2 = d.NullableString2
	}
	// embedded struct
	s.NullableString3 = NullStringMod{d.NullableString3}
	return s
}
