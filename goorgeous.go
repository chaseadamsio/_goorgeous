package goorgeous

func OrgToHTML(input string) string {
	return input
}

// Helpers
func skipChar(data []byte, start int, char byte) int {
	i := start
	for i < len(data) && charMatches(data[i], char) {
		i++
	}
	return i
}

func isSpace(char byte) bool {
	return charMatches(char, ' ')
}

func isEmpty(data []byte) bool {
	if len(data) == 0 {
		return true
	}

	for i := 0; i < len(data) && !charMatches(data[i], '\n'); i++ {
		if !charMatches(data[i], ' ') && !charMatches(data[i], '\t') {
			return false
		}
	}
	return true
}

func charMatches(a byte, b byte) bool {
	return a == b
}
