package utils

import "database/sql"

// NewNull creates a sql.Null[T] from an pointer.
func NewNull[T any](value *T) sql.Null[T] {
	var null sql.Null[T]

	if value != nil {
		null.Valid = true
		null.V = *value
	} else {
		null.Valid = false
	}

	return null
}

// PtrFromNull creates a pointer from sql.Null[T].
func PtrFromNull[T any](v sql.Null[T]) *T {
	if !v.Valid {
		return nil
	}

	vv := v.V
	return &vv
}
