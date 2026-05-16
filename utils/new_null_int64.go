package utils

import "database/sql"

// NewNullInt64 creates a sql.NullInt64 from an int64 pointer.
func NewNullInt64(value *int64) sql.NullInt64 {
	if value != nil {
		return sql.NullInt64{Valid: true, Int64: *value}
	}
	return sql.NullInt64{Valid: false}
}

// PtrFromNullInt64 creates an int64 pointer from sql.NullInt64.
func PtrFromNullInt64(v sql.NullInt64) *int64 {
	if !v.Valid {
		return nil
	}

	vv := v.Int64
	return &vv
}

