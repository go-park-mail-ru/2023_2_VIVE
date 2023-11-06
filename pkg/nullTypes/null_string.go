package nullTypes

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	Parent *sql.NullString
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Parent.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Parent.String)
}

// UnmarshalJSON for NullString
func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.Parent.String)
	ns.Parent.Valid = (err == nil)
	return err
}

// Creates new NullString
func NewNullString(value string, valid bool) NullString {
	return NullString{&sql.NullString{String: value, Valid: valid}}
}
