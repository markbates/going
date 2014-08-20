package nulls

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullString replaces sql.NullString with an implementation
// that supports proper JSON encoding/decoding.
type NullString sql.NullString

// NewNullString returns a new, properly instantiated
// NullString object.
func NewNullString(s string) NullString {
	return NullString{String: s, Valid: true}
}

// Scan implements the Scanner interface.
func (ns *NullString) Scan(value interface{}) error {
	n := sql.NullString{String: ns.String}
	err := n.Scan(value)
	ns.String, ns.Valid = n.String, n.Valid
	return err
}

// Value implements the driver Valuer interface.
func (ns NullString) Value() (driver.Value, error) {
	ns.Valid = ns.String != ""
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the propert representation of that value.
func (ns *NullString) UnmarshalJSON(text []byte) error {
	ns.Valid = false
	if string(text) == "null" {
		return nil
	}
	s := ""
	err := json.Unmarshal(text, &s)
	if err == nil {
		ns.String = s
		ns.Valid = true
	}
	return err
}
