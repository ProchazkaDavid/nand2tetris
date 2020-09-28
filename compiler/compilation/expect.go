package compilation

import (
	"errors"
	"fmt"

	"github.com/ProchazkaDavid/nand2tetris/compiler/token"
)

var (
	errExpectedIdentifier  = errors.New("expected identifier")
	errExpectedSymbol      = errors.New("expected symbol")
	errExpectedType        = errors.New("expected type")
	errExpectedKeyword     = errors.New("expected keyword")
	errUnexpectedType      = errors.New("unexpected type")
	errUnexpectedKeyword   = errors.New("unexpected keyword")
	errUnexpectedVMSegment = errors.New("unexpected VM segment")
)

func (e *Engine) expectOneOfTokens(tokens ...token.Type) {
	tt := e.tokenizer.TokenType()
	for _, t := range tokens {
		if t == tt {
			return
		}
	}

	e.handleError(fmt.Errorf("expected one of %v token type, got %v", tokens, tt))
}

func (e *Engine) expectIdentifier() {
	if e.tokenizer.TokenType() != token.Identifier {
		e.handleError(errExpectedIdentifier)
	}
}

func (e *Engine) isOneOfKeywords(keywords ...token.KeywordType) bool {
	if e.tokenizer.TokenType() != token.Keyword {
		return false
	}

	for _, k := range keywords {
		if k == e.tokenizer.Keyword() {
			return true
		}
	}

	return false
}

func (e *Engine) expectOneOfKeywords(keywords ...token.KeywordType) {
	if !e.isOneOfKeywords(keywords...) {
		e.handleError(errExpectedKeyword)
	}

	for _, k := range keywords {
		if k == e.tokenizer.Keyword() {
			return
		}
	}

	e.handleError(fmt.Errorf("expected one of %v keywords, got %v", keywords, e.tokenizer.Keyword()))
}

func (e *Engine) isCurrentSymbol(symbol string) bool {
	return e.tokenizer.TokenType() == token.Symbol && e.tokenizer.Symbol() == symbol
}

func (e *Engine) isCurrentKeyword(keyword token.KeywordType) bool {
	return e.tokenizer.TokenType() == token.Keyword && e.tokenizer.Keyword() == keyword
}

func (e *Engine) expectOneOfSymbols(symbols ...string) {
	if e.tokenizer.TokenType() != token.Symbol {
		e.handleError(errExpectedSymbol)
	}

	for _, s := range symbols {
		if s == e.tokenizer.Symbol() {
			return
		}
	}

	e.handleError(fmt.Errorf("expected one of %v symbols, got %v", symbols, e.tokenizer.Symbol()))
}

func (e *Engine) expectType() {
	switch e.tokenizer.TokenType() {
	case token.Keyword:
		e.expectOneOfKeywords(token.Int, token.Char, token.Boolean)
	case token.Identifier:
		// Class type
	default:
		e.handleError(errUnexpectedType)
	}
}
