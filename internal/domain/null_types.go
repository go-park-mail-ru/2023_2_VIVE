package domain

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = (err == nil)

	return err
}

/*type NullDate struct {
	sql.NullTime
}

func (n *NullDate) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	year, month, day := n.Time.Date()
	val := fmt.Sprintf("%d-%d-%d", year, int(month), day)

	return []byte(val), nil
}

func (n *NullDate) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &n.Time)
	n.Valid = (err == nil)

	return err
}*/
