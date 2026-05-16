package utils

// UintOrZero func for return uint value or zero if pointer is nil.
func UintOrZero(value *uint) uint {
	if value == nil {
		return 0
	}

	return *value
}
