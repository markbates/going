package nulls

import (
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
)

// NullByteSlice adds an implementation for []byte
// that supports proper JSON encoding/decoding.
type NullByteSlice struct {
	ByteSlice []byte
	Valid     bool // Valid is true if NullByteSlice is not NULL
}

// NewNullByteSlice returns a new, properly instantiated
// NullByteSlice object.
func NewNullByteSlice(b []byte) NullByteSlice {
	return NullByteSlice{ByteSlice: b, Valid: true}
}

// Scan implements the Scanner interface.
func (ns *NullByteSlice) Scan(value interface{}) error {
	n := sql.NullString{String: base64.StdEncoding.EncodeToString(ns.ByteSlice)}
	err := n.Scan(value)
	//ns.Float32, ns.Valid = float32(n.Float64), n.Valid
	ns.ByteSlice, err = base64.StdEncoding.DecodeString(n.String)
	ns.Valid = n.Valid
	return err
}

// Value implements the driver Valuer interface.
func (ns NullByteSlice) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return base64.StdEncoding.EncodeToString(ns.ByteSlice), nil
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns NullByteSlice) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.ByteSlice)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the propert representation of that value.
func (ns *NullByteSlice) UnmarshalJSON(text []byte) error {
	ns.Valid = false
	if string(text) == "null" {
		return nil
	}

	ns.ByteSlice = text
	ns.Valid = true
	return nil
}
