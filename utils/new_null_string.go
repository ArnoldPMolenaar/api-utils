package utils

import "database/sql"

// NewNullString creates a sql.NullString from a string pointer.
func NewNullString(value *string) sql.NullString {
	if value != nil && *value != "" {
		return sql.NullString{Valid: true, String: *value}
	}
	return sql.NullString{Valid: false}
}

// PtrFromNullString creates a string pointer from sql.NullString.
func PtrFromNullString(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}

	vv := v.String
	return &vv
}
