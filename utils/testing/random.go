package testing

import "math/rand"

const (
	Alpha        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numeric      = "0123456789"
	Punctuation  = "!`~@#$%^&*()_+-=[];',./{}:\"<>?\\|"
	VisibleAscii = Alpha + Numeric + Punctuation
	Unicode      = "\033\u2500\u2502\u2501"
)

func GenerateRandomString(length int, validChars string) string {
	var result []byte = make([]byte, length)
	for bi := range result {
		index := rand.Intn(len(validChars))
		result[bi] = validChars[index]
	}
	return string(result)
}

func GenerateRandomAlnumString(length int) string {
	return GenerateRandomString(length, Alpha+Numeric)
}

func GenerateRandomAsciiString(length int) string {
	return GenerateRandomString(length, VisibleAscii)
}

func GenerateRandomUnicodeString(length int) string {
	return GenerateRandomString(length, Unicode)
}
