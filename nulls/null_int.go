package nulls

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

// NullInt adds an implementation for int
// that supports proper JSON encoding/decoding.
type NullInt struct {
	Int   int
	Valid bool // Valid is true if Int is not NULL
}

// NewNullInt returns a new, properly instantiated
// NullInt object.
func NewNullInt(i int) NullInt {
	return NullInt{Int: i, Valid: true}
}

// Scan implements the Scanner interface.
func (ns *NullInt) Scan(value interface{}) error {
	n := sql.NullInt64{Int64: int64(ns.Int)}
	err := n.Scan(value)
	ns.Int, ns.Valid = int(n.Int64), n.Valid
	return err
}

// Value implements the driver Valuer interface.
func (ns NullInt) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return int64(ns.Int), nil
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns NullInt) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Int)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the propert representation of that value.
func (ns *NullInt) UnmarshalJSON(text []byte) error {
	txt := string(text)
	ns.Valid = true
	if txt == "null" {
		ns.Valid = false
		return nil
	}
	i, err := strconv.ParseInt(txt, 10, strconv.IntSize)
	if err != nil {
		ns.Valid = false
		return err
	}
	j := int(i)
	ns.Int = j
	return nil
}
