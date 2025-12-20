package util

import "strings"

const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// GenerateSortKey generates a sort key between two keys
// If prev is empty, generates key before next
// If next is empty, generates key after prev
// If both empty, returns middle key
func GenerateSortKey(prev, next string) string {
	if prev == "" && next == "" {
		return "U" // Middle of alphabet
	}

	if prev == "" {
		// Insert before next
		return decrementKey(next)
	}

	if next == "" {
		// Insert after prev
		return incrementKey(prev)
	}

	// Insert between prev and next
	return midKey(prev, next)
}

func incrementKey(key string) string {
	if key == "" {
		return "U"
	}

	runes := []rune(key)
	for i := len(runes) - 1; i >= 0; i-- {
		idx := strings.IndexRune(chars, runes[i])
		if idx < len(chars)-1 {
			runes[i] = rune(chars[idx+1])
			return string(runes)
		}
		runes[i] = '0'
	}
	return string(runes) + "0"
}

func decrementKey(key string) string {
	if key == "" {
		return "U"
	}

	runes := []rune(key)
	for i := len(runes) - 1; i >= 0; i-- {
		idx := strings.IndexRune(chars, runes[i])
		if idx > 0 {
			runes[i] = rune(chars[idx-1])
			return string(runes)
		}
	}
	return "0" + key
}

func midKey(prev, next string) string {
	// Pad to same length
	maxLen := len(prev)
	if len(next) > maxLen {
		maxLen = len(next)
	}

	for len(prev) < maxLen {
		prev += "0"
	}
	for len(next) < maxLen {
		next += "0"
	}

	result := make([]byte, maxLen)
	for i := 0; i < maxLen; i++ {
		p := strings.IndexByte(chars, prev[i])
		n := strings.IndexByte(chars, next[i])
		mid := (p + n) / 2
		result[i] = chars[mid]
	}

	res := string(result)
	if res == prev {
		return prev + "U"
	}
	return res
}
