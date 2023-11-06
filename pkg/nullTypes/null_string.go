package nullTypes

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON for NullString
func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = (err == nil)
	return err
}

// Retruns only stored value or nil
func (ni *NullString) GetValue() driver.Value {
	res, _ := ni.Value()
	return res
}

// Creates new NullString
func NewNullString(value string, valid bool) NullString {
	var res NullString
	if valid {
		res.String = value
		res.Valid = valid
	} else {
		res.Valid = valid
	}
	return res
}
