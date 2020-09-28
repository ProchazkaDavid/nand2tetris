package token

// Symbols represents group of language symbols
var Symbols = [...]string{
	"{",
	"}",
	"(",
	")",
	"[",
	"]",
	".",
	",",
	";",
	"+",
	"-",
	"*",
	"/",
	"&",
	"|",
	"<",
	">",
	"=",
	"~",
}

// IsSymbol checks if the input is symbol
func IsSymbol(input string) bool {
	for _, s := range Symbols {
		if input == s {
			return true
		}
	}

	return false
}

// ExpressionSymbols represents group of expression symbols
var ExpressionSymbols = [...]string{
	"+",
	"-",
	"*",
	"/",
	"&",
	"|",
	"<",
	">",
	"=",
}

// IsExpressionSymbol checks if the input is expression symbol
func IsExpressionSymbol(input string) bool {
	for _, s := range ExpressionSymbols {
		if input == s {
			return true
		}
	}

	return false
}
