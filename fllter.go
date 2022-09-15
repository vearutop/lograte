package main

// filterAlphanumeric replaces a-zA-Z_-% sequences that have at least one digit with X.
// Does not allocate, uses original slice.
func filterAlphanumeric(data []byte, l int) []byte {
	hasDigit := false
	wordStart := -1

	res := data[:0]

	var i int

	for i = 0; i < len(data); i++ {
		c := data[i]

		isAlpha := false

		switch {
		case c >= 'a' && c <= 'z':
			isAlpha = true
		case c >= 'A' && c <= 'Z':
			isAlpha = true
		case c >= '0' && c <= '9':
			isAlpha = true
			hasDigit = true
		case c == '-', c == '_', c == '%':
			isAlpha = true
		}

		// Finish current word.
		if wordStart >= 0 && !isAlpha {
			if hasDigit {
				res = append(res, 'X')
			} else {
				res = append(res, data[wordStart:i]...)
			}

			wordStart = -1
			hasDigit = false
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
