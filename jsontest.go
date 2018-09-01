package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
)

type NullString struct {
	sql.NullString
}

func (s *NullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}
	return json.Marshal("")
}

type TestStruct struct {
	TypeInt            int        `json:"typeInt"`
	TypeString         string     `json:"typeString"`
	TypeStringPtr      *string    `json:"typeStringPtr"`
	TypeNullableString NullString `json:"typeNullableString"`
}

func TestEncode() {
	data := &TestStruct{
		TypeInt:       1,
		TypeString:    "HogeFuga",
		TypeStringPtr: nil,
		TypeNullableString: NullString{sql.NullString{
			String: "N/C",
			Valid:  false,
		}},
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(data)

	fmt.Println(buf.String())
}

func main() {
	TestEncode()
}
