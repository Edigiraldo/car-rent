package utils

// Checks if a given value is in a given slice
func IsInSlice[T comparable](s []T, value T) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}

	return false
}
