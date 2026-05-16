package utils

import "database/sql"

// NewNullInt16 creates a sql.NullInt16 from a int16 pointer.
func NewNullInt16(value *int16) sql.NullInt16 {
	if value != nil {
		return sql.NullInt16{Valid: true, Int16: *value}
	}
	return sql.NullInt16{Valid: false}
}

// PtrFromNullInt16 creates an int16 pointer from sql.NullInt16.
func PtrFromNullInt16(v sql.NullInt16) *int16 {
	if !v.Valid {
		return nil
	}

	vv := v.Int16
	return &vv
}
