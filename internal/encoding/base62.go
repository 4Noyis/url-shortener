package encoding

import "strings"

const BASE62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func EncodeIntToBase62(n int64) string {
	if n == 0 {
		return string(BASE62[0])
	}

	if n < 0 {
		return ""
	}

	base := int64(len(BASE62))
	var result strings.Builder

	for n > 0 {
		rem := n % base
		result.WriteByte(BASE62[rem])
		n = n / base
	}

	return reverse(result.String())
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}