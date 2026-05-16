package utils

import "database/sql"

// NewNullInt32 creates a sql.NullInt32 from an int32 pointer.
func NewNullInt32(value *int32) sql.NullInt32 {
	if value != nil {
		return sql.NullInt32{Valid: true, Int32: *value}
	}
	return sql.NullInt32{Valid: false}
}

// PtrFromNullInt32 creates an int32 pointer from sql.NullInt32.
func PtrFromNullInt32(v sql.NullInt32) *int32 {
	if !v.Valid {
		return nil
	}

	vv := v.Int32
	return &vv
}
