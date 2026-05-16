package utils

import "database/sql"

// NewNullBool creates a sql.NullBool from a bool pointer.
func NewNullBool(value *bool) sql.NullBool {
	if value != nil {
		return sql.NullBool{Valid: true, Bool: *value}
	}
	return sql.NullBool{Valid: false}
}

// PtrFromNullBool creates a bool pointer from sql.NullBool.
func PtrFromNullBool(v sql.NullBool) *bool {
	if !v.Valid {
		return nil
	}

	vv := v.Bool
	return &vv
}
