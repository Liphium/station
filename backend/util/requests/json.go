package requests

// Type alias for map to make it more easily accessible
type Map = map[string]interface{}

// Get a boolean or other in case it's not there
func ValueOr[T any](m Map, key string, other T) T {
	val, valid := m[key]
	if !valid {
		return other
	}
	switch val := val.(type) {
	case T:
		return val
	}
	return other
}
