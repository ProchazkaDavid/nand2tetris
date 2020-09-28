package token

// Type represents token type
type Type int

// Token types
const (
	Keyword Type = iota
	Symbol
	Identifier
	IntegerConstant
	StringConstant
)
