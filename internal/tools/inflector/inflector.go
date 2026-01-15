package inflector

import (
	"strings"
	"unicode"
)

func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	str := []rune(s)
	return string(unicode.ToUpper(str[0])) + string(str[1:])
}

func FormatSentence(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	str := Capitalize(s)
	last := s[len(s)-1:]

	if last != "." && last != "?" && last != "!" {
		return str + "."
	}

	return str
}
