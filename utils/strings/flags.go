package strings

func IsNonEmpty(data []byte) bool {
	for _, datum := range data {
		if datum > 32 && datum < 127 {
			return true
		}
	}
	return false
}
