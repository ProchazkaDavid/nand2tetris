package compilation

import (
	"fmt"

	"github.com/ProchazkaDavid/nand2tetris/compiler/symbol"
	"github.com/ProchazkaDavid/nand2tetris/compiler/token"
)

// CompileClass compiles a complete class.
func (e *Engine) CompileClass() {
	e.advance()
	e.expectOneOfKeywords(token.Class)

	e.advance()
	e.expectIdentifier()

	e.className = e.tokenizer.Identifier()

	e.advance()
	e.expectOneOfSymbols("{")

	e.advance()

	variables := 0
	if e.isOneOfKeywords(token.Static, token.Field) {
		variables = e.compileClassVariableDeclarations()
	}

	if e.isOneOfKeywords(token.Constructor, token.Function, token.Method) {
		e.compileSubroutineDeclarations(variables)
	}

	e.expectOneOfSymbols("}")
}

// compileClassVariableDeclarations compiles a static variable declaration,
// or a field declaration.
func (e *Engine) compileClassVariableDeclarations() (variables int) {
	for e.isOneOfKeywords(token.Static, token.Field) {
		kind := symbol.Static
		if e.isCurrentKeyword(token.Field) {
			kind = symbol.Field
			variables++
		}

		e.advance()
		e.expectType()

		variableType := e.getVariableType()

		e.advance()
		e.expectIdentifier()

		e.symbolTable.Define(e.tokenizer.Identifier(), variableType, kind)

		e.advance()
		for ; e.isCurrentSymbol(","); e.advance() {
			e.advance()
			e.expectIdentifier()

			e.symbolTable.Define(e.tokenizer.Identifier(), variableType, kind)
			if kind == symbol.Field {
				variables++
			}
		}

		e.expectOneOfSymbols(";")
		e.advance()
	}

	return variables
}

// compileSubroutineDeclarations compiles a complete method, function,
// or constructor.
func (e *Engine) compileSubroutineDeclarations(classVariables int) {
	for e.isOneOfKeywords(token.Constructor, token.Function, token.Method) {
		e.symbolTable.NewSubroutine()
		e.ifCounter = 0
		e.whileCounter = 0

		subroutineType := e.tokenizer.Keyword()

		e.advance()
		if !e.isOneOfKeywords(token.Void) {
			e.expectType()
		}

		e.advance()
		e.expectIdentifier()

		function := e.tokenizer.Identifier()

		e.advance()
		e.expectOneOfSymbols("(")

		if subroutineType == token.Method {
			e.symbolTable.Define("this", e.className, symbol.Arg)
		}

		e.advance()
		e.compileParameterList()

		e.expectOneOfSymbols(")")

		e.advance()
		e.compileSubroutineBody(fmt.Sprintf("%s.%s", e.className, function), subroutineType, classVariables)
	}
}
