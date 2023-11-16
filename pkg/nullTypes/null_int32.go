package nullTypes

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type NullInt32 struct {
	sql.NullInt32
}

// MarshalJSON for NullInt
func (ni *NullInt32) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int32)
}

// UnmarshalJSON for NullInt
func (ni *NullInt32) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Int32)
	ni.Valid = (err == nil)
	return err
}

// Retruns only stored value or nil
func (ni *NullInt32) GetValue() driver.Value {
	res, _ := ni.Value()
	return res
}

// Creates new NullInt
func NewNullInt(value int32, valid bool) NullInt32 {
	var res NullInt32
	if valid {
		res.Int32 = value
		res.Valid = valid
	} else {
		res.Valid = valid
	}
	return res
}
