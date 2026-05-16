package utils

import (
	"database/sql"
	"time"
)

// NewNullTime creates a sql.NullTime from a time pointer.
func NewNullTime(value *time.Time) sql.NullTime {
	if value != nil {
		return sql.NullTime{Valid: true, Time: *value}
	}
	return sql.NullTime{Valid: false}
}

// PtrFromNullTime creates a time pointer from sql.NullTime.
func PtrFromNullTime(v sql.NullTime) *time.Time {
	if !v.Valid {
		return nil
	}

	vv := v.Time
	return &vv
}
