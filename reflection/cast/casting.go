package cast

// As creates new slice instance and casts from elements to result slice.
// It will return true when from slice is empty
func As[R any, T any](from []T) ([]R, bool) {
	count := len(from)
	result := make([]R, count)

	for i := 0; i < count; i++ {
		switch t := any(from[i]).(type) {
		case R:
			result[i] = R(t)
		default:
			return nil, false
		}
	}

	return result, true
}
