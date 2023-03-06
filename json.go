package json_utils

import (
	"strings"
)

// Delete comments and BOM (if any) from JSON
func SanitizeJson(b []byte) []byte {
	b = RemoveBom(b)
	b = stripComments(b)
	return b
}


// States
const (
	emptyState = iota
	stringState
	singleLineCommentState
	multilineCommentState
)

// Remove comments with LL(2) lexer
func stripComments(b []byte) []byte {
	state := emptyState
	offset := 0
	var buffer strings.Builder
	var result strings.Builder
	commaIndex := -1

	for index, n := 0, len(b); index < n; index++ {
		ch := getAt(b, index)
		nextCh := getAt(b, index + 1)
		
		if checkState(state, emptyState) {
			// select next state
			if ch == '"' {
				// string
				state = stringState
				if commaIndex != -1 {
					commaIndex = -1
				}
			} else if ch == '/' {
				// may be a comment
				if nextCh == '/' {
					// single line comment
					buffer.Write(b[offset:index])
					offset = index
					state = singleLineCommentState
					index++;
				} else if nextCh == '*' {
					// multiline comment
					buffer.Write(b[offset:index])
					offset = index
					state = multilineCommentState
					index++
				}
			} else if commaIndex != -1 {
				if ch == '}' || ch == ']' {
					// deleting unnecessary commas
					buffer.Write(b[offset:index])
					result.Write([]byte(buffer.String())[1:])
					buffer.Reset()
					offset = index
					commaIndex = -1
				} else if ch != ' ' && ch != '\t' && ch != '\r' && ch != '\n' {
					// conditionally non-whitespace character after the comma
					buffer.Write(b[offset:index])
					offset = index
					commaIndex = -1
				}
			} else if ch == ',' {
				buffer.Write(b[offset:index])
				result.Write([]byte(buffer.String()))
				buffer.Reset()
				offset = index
				commaIndex = index
			}
		} else if checkState(state, stringState) {
			// string could be closed only by non shielded "
			if ch == '"' && !isEscaped(b, index) {
				state = emptyState
			}
		} else if checkState(state, singleLineCommentState) {
			// is single line comment could be finished
			if ch == '\r' {
				// finishing by \r(\n)?
				offset = index
				state = emptyState
				if nextCh == '\n' {
					index++
				}
			} else if ch == '\n' {
				// finishing by \n
				state = emptyState
				offset = index
			}
		} else if checkState(state, multilineCommentState) {
			// multiline comment finished only by paired */
			if ch == '*' && nextCh == '/' {
				index++
				state = emptyState
				offset = index + 1
			}
		}
	}

	if checkState(state, emptyState) {
		buffer.Write(b[offset:])
	}
	result.Write([]byte(buffer.String()))
	return []byte(result.String())
}

// Get byte value in given position.
func getAt(b []byte, pos int) (r byte) {
	if pos < 0 || pos >= len(b) {
		return
	}
	r = b[pos]
	return
}

// Check whether symbol is shielded.
func isEscaped(b []byte, pos int) bool {
	backslashCount := 0
	for pos--; pos >= 0 && b[pos] == '\\'; pos-- {
		backslashCount++
	}
	return backslashCount % 2 != 0
}

// Check state.
func checkState(state, mode int) bool {
	return state == mode
}

