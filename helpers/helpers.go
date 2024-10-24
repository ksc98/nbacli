package helpers

var TOGGLE_MAP = map[string]bool{}

func Toggle(key string) bool {
	if _, ok := TOGGLE_MAP[key]; ok {
		// delete(TOGGLE_MAP, key)
		return true
	} else {
		TOGGLE_MAP[key] = true
		return false
	}
}
