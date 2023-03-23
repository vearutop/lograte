package main

// filterDynamic replaces a-zA-Z_-% sequences that have at least one digit or 5+ consecutive consolants/vowels with X.
// Does not allocate, uses original slice.
func filterDynamic(data []byte, l int) []byte {
	hasDigit := false
	wordStart := -1
	maxConsecutive := 0
	consecutive := 0

	res := data[:0]

	var (
		i            int
		prevCharType byte
	)

	for i = 0; i < len(data); i++ {
		c := data[i]

		isAlpha := false
		var charType byte

		switch {
		case c >= 'a' && c <= 'z':
			if c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u' || c == 'y' || c == 'w' {
				charType = 'v' // Vowel.
			} else {
				charType = 'c' // Consonant.
			}

			isAlpha = true
		case c >= 'A' && c <= 'Z':
			if c == 'A' || c == 'E' || c == 'I' || c == 'O' || c == 'U' || c == 'Y' || c == 'W' {
				charType = 'v'
			} else {
				charType = 'c'
			}

			isAlpha = true
		case c >= '0' && c <= '9':
			isAlpha = true
			hasDigit = true
		case c == '-', c == '_', c == '%':
			isAlpha = true
		}

		if charType == prevCharType {
			consecutive++
		} else {
			if consecutive > maxConsecutive {
				maxConsecutive = consecutive
			}

			prevCharType = charType
			consecutive = 1
		}

		// Finish current word.
		if wordStart >= 0 && !isAlpha {
			if hasDigit || maxConsecutive > 4 {
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
		if hasDigit {
			res = append(res, 'X')
		} else {
			res = append(res, data[wordStart:i]...)
		}
	}

	return res
}
