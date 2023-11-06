package nullTypes

import (
	"database/sql"
	"encoding/json"
)

type NullInt struct {
	Parent sql.NullInt32
}

// MarshalJSON for NullInt
func (ni *NullInt) MarshalJSON() ([]byte, error) {
	if !ni.Parent.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Parent.Int32)
}

// UnmarshalJSON for NullInt
func (ni *NullInt) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Parent.Int32)
	ni.Parent.Valid = (err == nil)
	return err
}

// Creates new NullInt
func NewNullInt(value int, valid bool) NullInt {
	return NullInt{sql.NullInt32{Int32: int32(value), Valid: valid}}
}
