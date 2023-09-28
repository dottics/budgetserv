package budget

// NotEqualError returns true if the two errors are not equal.
func NotEqualError(e1, e2 error) bool {
	switch {
	case e1 == nil && e2 == nil:
		return false
	case e1 == nil || e2 == nil:
		return true
	case e1.Error() != e2.Error():
		return true
	default:
		return false
	}
}
