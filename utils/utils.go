package utils

// make ternary operator utility function for golang
func TernaryIF(condition bool, a, b interface{}) *interface{} {
	if condition {
		return &a
	}
	return &b
}
