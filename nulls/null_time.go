package nulls

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// NullTime replaces sql.NullTime with an implementation
// that supports proper JSON encoding/decoding.
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// NewNullTime returns a new, properly instantiated
// NullTime object.
func NewNullTime(t time.Time) NullTime {
	return NullTime{Time: t, Valid: true}
}

// Scan implements the Scanner interface.
func (ns *NullTime) Scan(value interface{}) error {
	ns.Time, ns.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullTime) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Time, nil
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns NullTime) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Time)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the propert representation of that value.
func (ns *NullTime) UnmarshalJSON(text []byte) error {
	ns.Valid = false
	txt := string(text)
	if txt == "null" || txt == "" {
		return nil
	}

	t := time.Time{}
	err := t.UnmarshalJSON(text)
	if err == nil {
		ns.Time = t
		ns.Valid = true
	}

	return err
}
