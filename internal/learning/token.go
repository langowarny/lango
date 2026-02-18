package learning

import "unicode/utf8"

const (
	asciiCharsPerToken = 4
	cjkRunesPerToken   = 2
	perMessageOverhead = 4
)

func estimateTokens(text string) int {
	var asciiCount, cjkCount int
	for i := 0; i < len(text); {
		r, size := utf8.DecodeRuneInString(text[i:])
		if isCJKRune(r) {
			cjkCount++
		} else {
			asciiCount++
		}
		i += size
	}
	return asciiCount/asciiCharsPerToken + cjkCount/cjkRunesPerToken
}

func isCJKRune(r rune) bool {
	if r >= 0x4E00 && r <= 0x9FFF {
		return true
	}
	if r >= 0x3400 && r <= 0x4DBF {
		return true
	}
	if r >= 0xAC00 && r <= 0xD7AF {
		return true
	}
	return false
}
