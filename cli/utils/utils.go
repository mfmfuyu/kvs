package utils

import "strings"

func Parse(text string) []string {
	var result = []string{}

	var quote rune
	var builder strings.Builder
	var escape = false

	for _, r := range text {
		if escape {
			builder.WriteRune(r)
			escape = false
			continue
		}

		if r == '\\' {
			escape = true
			continue
		}

		if quote != 0 {
			if r == quote {
				quote = 0
			} else {
				builder.WriteRune(r)
			}
			continue
		}

		switch r {
		case '\'', '"':
			quote = r
		case ' ':
			if builder.Len() > 0 {
				result = append(result, builder.String())
				builder.Reset()
			}
		default:
			builder.WriteRune(r)
		}
	}

	if quote != 0 {
		return nil
	}

	if builder.Len() > 0 {
		result = append(result, builder.String())
	}

	return result
}

func Digits(n int) int {
	if n == 0 {
		return 1
	}

	if n < 0 {
		n = -n
	}

	d := 0
	for n > 0 {
		n /= 10
		d++
	}

	return d
}
