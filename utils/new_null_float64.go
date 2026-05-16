package utils

import "database/sql"

// NewNullFloat64 creates a sql.NullFloat64 from a float64 pointer.
func NewNullFloat64(value *float64) sql.NullFloat64 {
	if value != nil {
		return sql.NullFloat64{Valid: true, Float64: *value}
	}
	return sql.NullFloat64{Valid: false}
}

// PtrFromNullFloat64 creates a float64 pointer from sql.NullFloat64.
func PtrFromNullFloat64(v sql.NullFloat64) *float64 {
	if !v.Valid {
		return nil
	}

	vv := v.Float64
	return &vv
}
