package editor

import "unicode"

type TokenType int

const (
	TokenPlain TokenType = iota
	TokenKeyword
	TokenString
	TokenNumber
	TokenComment
	TokenFunction
	TokenSpace
	TokenTab
	TokenType_
)

type Token struct {
	Type  TokenType
	Value []rune
	Start int
	End   int
}

type SyntaxRule struct {
	Keywords              []string
	Types                 []string
	LineComment           string
	MultiLineCommentStart string
	MultiLineCommentEnd   string
	StringDelimiters      []rune
}

var languageRules = map[string]SyntaxRule{
	"Go": {
		Keywords: []string{
			"break", "default", "func", "interface", "select",
			"case", "defer", "go", "map", "struct",
			"chan", "else", "goto", "package", "switch",
			"const", "fallthrough", "if", "range", "type",
			"continue", "for", "import", "return", "var",
		},
		Types: []string{
			"string", "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64", "bool", "byte", "rune",
			"error", "interface{}", "any",
		},
		LineComment:           "//",
		MultiLineCommentStart: "/*",
		MultiLineCommentEnd:   "*/",
		StringDelimiters:      []rune{'"', '`'},
	},
	"Python": {
		Keywords: []string{
			"False", "None", "True", "and", "as", "assert",
			"break", "class", "continue", "def", "del", "elif",
			"else", "except", "finally", "for", "from", "global",
			"if", "import", "in", "is", "lambda", "nonlocal",
			"not", "or", "pass", "raise", "return", "try",
			"while", "with", "yield",
		},
		Types: []string{
			"int", "float", "str", "list", "dict", "set",
			"tuple", "bool", "bytes", "object",
		},
		LineComment:           "#",
		MultiLineCommentStart: `"""`,
		MultiLineCommentEnd:   `"""`,
		StringDelimiters:      []rune{'"', '\''},
	},
	"JavaScript": {
		Keywords: []string{
			"break", "case", "catch", "class", "const", "continue",
			"debugger", "default", "delete", "do", "else", "export",
			"extends", "finally", "for", "function", "if", "import",
			"in", "instanceof", "new", "return", "super", "switch",
			"this", "throw", "try", "typeof", "var", "void", "while",
			"with", "yield", "let", "static", "enum", "await", "async",
		},
		Types: []string{
			"Array", "Boolean", "Date", "Error", "Function", "JSON",
			"Math", "Number", "Object", "RegExp", "String", "undefined",
			"null", "NaN", "Infinity", "Promise", "Map", "Set",
		},
		LineComment:           "//",
		MultiLineCommentStart: "/*",
		MultiLineCommentEnd:   "*/",
		StringDelimiters:      []rune{'"', '\'', '`'},
	},
	"Ruby": {
		Keywords: []string{
			"BEGIN", "END", "alias", "and", "begin", "break",
			"case", "class", "def", "defined?", "do", "else",
			"elsif", "end", "ensure", "false", "for", "if",
			"in", "module", "next", "nil", "not", "or", "redo",
			"rescue", "retry", "return", "self", "super", "then",
			"true", "undef", "unless", "until", "when", "while",
			"yield",
		},
		Types: []string{
			"Array", "Hash", "String", "Integer", "Float", "Symbol",
			"NilClass", "TrueClass", "FalseClass", "Numeric",
		},
		LineComment:           "#",
		MultiLineCommentStart: "=begin",
		MultiLineCommentEnd:   "=end",
		StringDelimiters:      []rune{'"', '\''},
	},
	"C": {
		Keywords: []string{
			"auto", "break", "case", "char", "const", "continue",
			"default", "do", "double", "else", "enum", "extern",
			"float", "for", "goto", "if", "inline", "int", "long",
			"register", "restrict", "return", "short", "signed",
			"sizeof", "static", "struct", "switch", "typedef",
			"union", "unsigned", "void", "volatile", "while",
		},
		Types: []string{
			"bool", "char", "double", "float", "int", "long",
			"short", "size_t", "void", "wchar_t", "unsigned",
			"signed", "uint8_t", "uint16_t", "uint32_t", "uint64_t",
			"int8_t", "int16_t", "int32_t", "int64_t",
		},
		LineComment:           "//",
		MultiLineCommentStart: "/*",
		MultiLineCommentEnd:   "*/",
		StringDelimiters:      []rune{'"'},
	},
	"C++": {
		Keywords: []string{
			"alignas", "alignof", "and", "and_eq", "asm", "auto",
			"bitand", "bitor", "bool", "break", "case", "catch",
			"char", "class", "compl", "const", "constexpr",
			"const_cast", "continue", "decltype", "default",
			"delete", "do", "double", "dynamic_cast", "else",
			"enum", "explicit", "export", "extern", "false",
			"float", "for", "friend", "goto", "if", "inline",
			"int", "long", "mutable", "namespace", "new",
			"noexcept", "not", "not_eq", "nullptr", "operator",
			"or", "or_eq", "private", "protected", "public",
			"register", "reinterpret_cast", "return", "short",
			"signed", "sizeof", "static", "static_assert",
			"static_cast", "struct", "switch", "template", "this",
			"thread_local", "throw", "true", "try", "typedef",
			"typeid", "typename", "union", "unsigned", "using",
			"virtual", "void", "volatile", "wchar_t", "while",
			"xor", "xor_eq",
		},
		Types: []string{
			"bool", "char", "char8_t", "char16_t", "char32_t",
			"double", "float", "int", "long", "short", "signed",
			"unsigned", "void", "wchar_t", "size_t", "string",
			"vector", "map", "set", "list", "queue", "stack",
			"array", "deque", "pair", "tuple",
		},
		LineComment:           "//",
		MultiLineCommentStart: "/*",
		MultiLineCommentEnd:   "*/",
		StringDelimiters:      []rune{'"'},
	},
	"Java": {
		Keywords: []string{
			"abstract", "assert", "break", "case", "catch",
			"class", "const", "continue", "default", "do",
			"else", "enum", "extends", "final", "finally",
			"for", "goto", "if", "implements", "import",
			"instanceof", "interface", "native", "new",
			"package", "private", "protected", "public",
			"return", "static", "strictfp", "super", "switch",
			"synchronized", "this", "throw", "throws", "transient",
			"try", "void", "volatile", "while",
		},
		Types: []string{
			"boolean", "byte", "char", "double", "float",
			"int", "long", "short", "String", "Object",
			"Integer", "Long", "Float", "Double", "Boolean",
			"Character", "Byte", "Short", "List", "Map",
			"Set", "Collection", "ArrayList", "HashMap",
			"HashSet", "Vector", "Array",
		},
		LineComment:           "//",
		MultiLineCommentStart: "/*",
		MultiLineCommentEnd:   "*/",
		StringDelimiters:      []rune{'"'},
	},
}

func tokenizeLine(line []rune, lang string, inComment bool) ([]Token, bool) {
	var tokens []Token
	rules, exists := languageRules[lang]
	if !exists {
		return []Token{{Type: TokenPlain, Value: line, Start: 0, End: len(line)}}, false
	}

	pos := 0
	for pos < len(line) {
		if inComment && len(rules.MultiLineCommentEnd) > 0 {
			if commentEndIndex := runesIndex(line[pos:], []rune(rules.MultiLineCommentEnd)); commentEndIndex >= 0 {
				tokens = append(tokens, Token{
					Type:  TokenComment,
					Value: line[pos : pos+commentEndIndex+len(rules.MultiLineCommentEnd)],
					Start: pos,
					End:   pos + commentEndIndex + len(rules.MultiLineCommentEnd),
				})
				pos += commentEndIndex + len(rules.MultiLineCommentEnd)
				inComment = false
				continue
			}
			tokens = append(tokens, Token{
				Type:  TokenComment,
				Value: line[pos:],
				Start: pos,
				End:   len(line),
			})
			return tokens, true
		}

		if unicode.IsSpace(line[pos]) {
			if line[pos] == ' ' {
				tokens = append(tokens, Token{
					Type:  TokenSpace,
					Value: line[pos : pos+1],
					Start: pos,
					End:   pos + 1,
				})
				pos++
				continue
			}
			if line[pos] == '\t' {
				tokens = append(tokens, Token{
					Type:  TokenTab,
					Value: line[pos : pos+1],
					Start: pos,
					End:   pos + 1,
				})
				pos++
				continue
			}
			pos++
			continue
		}

		if len(rules.LineComment) > 0 && hasPrefix(line[pos:], rules.LineComment) {
			tokens = append(tokens, Token{
				Type:  TokenComment,
				Value: line[pos:],
				Start: pos,
				End:   len(line),
			})
			break
		}

		if len(rules.MultiLineCommentStart) > 0 && hasPrefix(line[pos:], rules.MultiLineCommentStart) {
			commentEndIndex := runesIndex(line[pos+len(rules.MultiLineCommentStart):], []rune(rules.MultiLineCommentEnd))
			if commentEndIndex >= 0 {
				endPos := pos + len(rules.MultiLineCommentStart) + commentEndIndex + len(rules.MultiLineCommentEnd)
				tokens = append(tokens, Token{
					Type:  TokenComment,
					Value: line[pos:endPos],
					Start: pos,
					End:   endPos,
				})
				pos = endPos
				continue
			}
			tokens = append(tokens, Token{
				Type:  TokenComment,
				Value: line[pos:],
				Start: pos,
				End:   len(line),
			})
			return tokens, true
		}

		if isStringStart(line[pos], rules.StringDelimiters) {
			strEnd := findStringEnd(line[pos+1:], line[pos])
			if strEnd >= 0 {
				tokens = append(tokens, Token{
					Type:  TokenString,
					Value: line[pos : pos+strEnd+2],
					Start: pos,
					End:   pos + strEnd + 2,
				})
				pos += strEnd + 2
				continue
			}
		}

		if unicode.IsDigit(line[pos]) {
			numEnd := consumeNumber(line[pos:])
			tokens = append(tokens, Token{
				Type:  TokenNumber,
				Value: line[pos : pos+numEnd],
				Start: pos,
				End:   pos + numEnd,
			})
			pos += numEnd
			continue
		}

		if unicode.IsLetter(line[pos]) || line[pos] == '_' {
			wordEnd := consumeWord(line[pos:])
			word := string(line[pos : pos+wordEnd])

			tokenType := TokenPlain
			if isKeyword(word, rules.Keywords) {
				tokenType = TokenKeyword
			} else if isType(word, rules.Types) {
				tokenType = TokenType_
			} else if isFunctionCall(line[pos+wordEnd:]) {
				tokenType = TokenFunction
			}

			tokens = append(tokens, Token{
				Type:  tokenType,
				Value: line[pos : pos+wordEnd],
				Start: pos,
				End:   pos + wordEnd,
			})
			pos += wordEnd
			continue
		}

		tokens = append(tokens, Token{
			Type:  TokenPlain,
			Value: line[pos : pos+1],
			Start: pos,
			End:   pos + 1,
		})
		pos++
	}

	return tokens, inComment
}

func hasPrefix(text []rune, prefix string) bool {
	if len(text) < len(prefix) {
		return false
	}
	for i, r := range prefix {
		if text[i] != r {
			return false
		}
	}
	return true
}

func runesIndex(text []rune, substr []rune) int {
	if len(substr) == 0 {
		return 0
	}
	if len(substr) > len(text) {
		return -1
	}
	for i := 0; i <= len(text)-len(substr); i++ {
		matched := true
		for j := 0; j < len(substr); j++ {
			if text[i+j] != substr[j] {
				matched = false
				break
			}
		}
		if matched {
			return i
		}
	}
	return -1
}

func isStringStart(r rune, delimiters []rune) bool {
	for _, d := range delimiters {
		if r == d {
			return true
		}
	}
	return false
}

func findStringEnd(text []rune, delimiter rune) int {
	for i := 0; i < len(text); i++ {
		if text[i] == '\\' {
			i++
			continue
		}
		if text[i] == delimiter {
			return i
		}
	}
	return -1
}

func consumeNumber(text []rune) int {
	i := 0
	foundDot := false
	for i < len(text) {
		if unicode.IsDigit(text[i]) {
			i++
			continue
		}
		if text[i] == '.' && !foundDot {
			foundDot = true
			i++
			continue
		}
		break
	}
	return i
}

func consumeWord(text []rune) int {
	i := 0
	for i < len(text) && (unicode.IsLetter(text[i]) || unicode.IsDigit(text[i]) || text[i] == '_') {
		i++
	}
	return i
}

func isKeyword(word string, keywords []string) bool {
	for _, k := range keywords {
		if word == k {
			return true
		}
	}
	return false
}

func isType(word string, types []string) bool {
	for _, t := range types {
		if word == t {
			return true
		}
	}
	return false
}

func isFunctionCall(text []rune) bool {
	for _, r := range text {
		if unicode.IsSpace(r) {
			continue
		}
		return r == '('
	}
	return false
}
