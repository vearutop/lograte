// Package filter provides a function to filter dynamic values/identifiers from a string.
package filter

// Dynamic replaces a-zA-Z_-% sequences that have at least one digit or consecutively 5+ consonants, 4+ vowels with X.
//
// This is useful for filtering out dynamic values from log lines.
// Filtering is best effort, it does not guarantee that all dynamic values are filtered.
// Does not allocate, mutates original slice.
func Dynamic(data []byte, l int) []byte {
	hasDigit := false
	wordStart := -1
	maxConsecutive := 0
	consecutive := 0
	res := data[:0]

	var (
		i            int
		prevCharType byte
		maxCharType  byte
	)

	for i = 0; i < len(data); i++ {
		c := data[i]

		isAlpha := false
		charType := byte(0)

		switch {
		case c >= 'a' && c <= 'z':
			if c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u' || c == 'y' {
				charType = 'v' // Vowel.
			} else {
				charType = 'c' // Consonant.
			}

			isAlpha = true
		case c >= 'A' && c <= 'Z':
			if c == 'A' || c == 'E' || c == 'I' || c == 'O' || c == 'U' || c == 'Y' {
				charType = 'v'
			} else {
				charType = 'c'
			}

			isAlpha = true
		case c >= '0' && c <= '9':
			isAlpha = true
			hasDigit = true
		case c == '_', c == '%', c == '-':
			isAlpha = true
		}

		if charType == prevCharType {
			consecutive++
		} else {
			if consecutive > maxConsecutive {
				maxConsecutive = consecutive
				maxCharType = prevCharType
			}

			prevCharType = charType
			consecutive = 1
		}

		// Finish current word.
		if wordStart >= 0 && !isAlpha {
			if hasDigit || (maxConsecutive > 3 && maxCharType == 'v') || (maxConsecutive > 4 && maxCharType == 'c') {
				res = append(res, 'X')
			} else {
				res = append(res, data[wordStart:i]...)
			}

			wordStart = -1
			hasDigit = false
			maxConsecutive = 0
			consecutive = 0
		}

		if wordStart == -1 {
			if isAlpha {
				// Starting new word.
				wordStart = i
			} else {
				// Adding current char.
				res = append(res, c)
			}
		}

		if l > 0 && len(res) >= l {
			return res
		}
	}

	if wordStart >= 0 {
		if hasDigit || (maxConsecutive > 3 && maxCharType == 'v') || (maxConsecutive > 4 && maxCharType == 'c') {
			res = append(res, 'X')
		} else {
			res = append(res, data[wordStart:i]...)
		}
	}

	return res
}
