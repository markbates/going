package nulls

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullBool replaces sql.NullBool with an implementation
// that supports proper JSON encoding/decoding.
type NullBool sql.NullBool

// NewNullBool returns a new, properly instantiated
// NullBoll object.
func NewNullBool(b bool) NullBool {
	return NullBool{Bool: b, Valid: true}
}

// Scan implements the Scanner interface.
func (ns *NullBool) Scan(value interface{}) error {
	n := sql.NullBool{Bool: ns.Bool}
	err := n.Scan(value)
	ns.Bool, ns.Valid = n.Bool, n.Valid
	return err
}

// Value implements the driver Valuer interface.
func (ns NullBool) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Bool, nil
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns NullBool) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Bool)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the proper representation of that value. The strings
// "true" and "t" will be considered "true", all other
// values will be considered "false". Null values will
//be set by setting Valid = false
func (ns *NullBool) UnmarshalJSON(text []byte) error {
	t := string(text)
	if t == "true" || t == "t" {
		ns.Valid = true
		ns.Bool = true
		return nil
	}
	ns.Bool = false
	ns.Valid = false
	return nil
}
