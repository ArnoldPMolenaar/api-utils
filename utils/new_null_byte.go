package utils

import "database/sql"

// NewNullByte creates a sql.NullByte from a byte pointer.
func NewNullByte(value *byte) sql.NullByte {
	if value != nil {
		return sql.NullByte{Valid: true, Byte: *value}
	}
	return sql.NullByte{Valid: false}
}

// PtrFromNullByte creates a byte pointer from sql.NullByte.
func PtrFromNullByte(v sql.NullByte) *byte {
	if !v.Valid {
		return nil
	}

	vv := v.Byte
	return &vv
}
