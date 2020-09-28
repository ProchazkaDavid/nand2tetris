package token

// KeywordType represents keyword type
type KeywordType string

// Keyword types
const (
	Unknown     KeywordType = ""
	Class       KeywordType = "class"
	Method      KeywordType = "method"
	Function    KeywordType = "function"
	Constructor KeywordType = "constructor"
	Int         KeywordType = "int"
	Boolean     KeywordType = "boolean"
	Char        KeywordType = "char"
	Void        KeywordType = "void"
	Var         KeywordType = "var"
	Static      KeywordType = "static"
	Field       KeywordType = "field"
	Let         KeywordType = "let"
	Do          KeywordType = "do"
	If          KeywordType = "if"
	Else        KeywordType = "else"
	While       KeywordType = "while"
	Return      KeywordType = "return"
	True        KeywordType = "true"
	False       KeywordType = "false"
	Null        KeywordType = "null"
	This        KeywordType = "this"
)

// Keywords represents mapping between string representation of Type and Type itself
var Keywords = [...]KeywordType{
	Class,
	Constructor,
	Function,
	Method,
	Field,
	Static,
	Var,
	Int,
	Char,
	Boolean,
	Void,
	True,
	False,
	Null,
	This,
	Let,
	Do,
	If,
	Else,
	While,
	Return,
}

// IsKeyword checks if the input is keyword
func IsKeyword(input string) bool {
	for _, keyword := range Keywords {
		if input == string(keyword) {
			return true
		}
	}

	return false
}
