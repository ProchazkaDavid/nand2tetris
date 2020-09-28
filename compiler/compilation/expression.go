package compilation

import (
	"fmt"

	"github.com/ProchazkaDavid/nand2tetris/compiler/token"
	"github.com/ProchazkaDavid/nand2tetris/compiler/vm"
)

// compileExpression compiles an expression.
func (e *Engine) compileExpression() {
	e.compileTerm()

	for token.IsExpressionSymbol(e.tokenizer.Symbol()) {
		operation := e.tokenizer.Symbol()

		e.advance()
		e.compileTerm()

		e.vm.WriteArithmetic(operation)
	}
}

// compileTerm compiles a term. If the current token is an identifier,
// the routine must distinguish between a variable, an array entry, or a
// subroutine call. A single look-ahead token, which may be one of "[", "(", or ".",
// suffices to distinguish between the possibilitios. Any other token is not part
// of this term and should not be advanced over.
func (e *Engine) compileTerm() {
	e.expectOneOfTokens(token.IntegerConstant, token.StringConstant, token.Keyword, token.Identifier, token.Symbol)

	switch e.tokenizer.TokenType() {
	case token.IntegerConstant:
		e.vm.WritePush(vm.Constant, e.tokenizer.IntValue())
		e.advance()

	case token.StringConstant:
		e.vm.WriteString(e.tokenizer.StringValue())
		e.advance()

	case token.Keyword:
		e.compileTermKeyword()

	case token.Identifier:
		e.compileTermIdentifier()

	case token.Symbol:
		e.compileSymbol()
	}
}

// compileTermKeyword compiles keywords True, False, Null, and This.
// The keywords push following values onto the stack:
//  True -> -1
//  False, Null -> 0
//  This -> pointer to the current object
func (e *Engine) compileTermKeyword() {
	e.expectOneOfKeywords(token.True, token.False, token.Null, token.This)

	var segment vm.Segment

	switch e.tokenizer.Keyword() {
	case token.True, token.False, token.Null:
		segment = vm.Constant
	case token.This:
		segment = vm.Pointer
	default:
		e.handleError(errUnexpectedKeyword)
	}

	e.vm.WritePush(segment, 0)

	if e.isCurrentKeyword(token.True) {
		e.vm.WriteArithmetic("~")
	}

	e.advance()
}

func (e *Engine) compileTermIdentifier() {
	termName := e.tokenizer.Identifier()

	e.advance()

	if e.isCurrentSymbol("(") {
		e.compileTermIdentifierFunction(termName)
		e.advance()
		return
	}

	termSegment := vm.GetSegment(e.symbolTable.KindOf(termName))
	symbolTableIndex := e.symbolTable.IndexOf(termName)

	switch {
	case e.isCurrentSymbol("["):
		e.compileTermIdentifierArray(termSegment, symbolTableIndex)
	case e.isCurrentSymbol("."):
		e.compileTermIdentifierMethod(termName, termSegment, symbolTableIndex)
	default:
		e.vm.WritePush(termSegment, symbolTableIndex)
		return
	}

	e.advance()
}

func (e *Engine) compileTermIdentifierArray(segment vm.Segment, symbolTableIndex int) {
	e.expectOneOfSymbols("[")

	e.advance()
	e.compileExpression()

	e.vm.WritePush(segment, symbolTableIndex)
	e.vm.WriteArithmetic("+")
	e.vm.WritePop(vm.Pointer, 1)
	e.vm.WritePush(vm.That, 0)

	e.expectOneOfSymbols("]")
}

func (e *Engine) compileTermIdentifierFunction(term string) {
	e.expectOneOfSymbols("(")
	e.advance()

	expressions := e.compileExpressionList()

	e.expectOneOfSymbols(")")
	e.vm.WriteCall(term, expressions)
}

func (e *Engine) compileTermIdentifierMethod(term string, segment vm.Segment, symbolTableIndex int) {
	e.expectOneOfSymbols(".")

	e.advance()
	e.expectIdentifier()

	function := e.tokenizer.Identifier()

	e.advance()
	e.expectOneOfSymbols("(")

	e.advance()
	expressions := e.compileExpressionList()

	class, ok := e.symbolTable.TypeOf(term)
	if !ok {
		class = term
	} else {
		e.vm.WritePush(segment, symbolTableIndex)
		expressions++
	}

	e.vm.WriteCall(fmt.Sprintf("%s.%s", class, function), expressions)

	e.expectOneOfSymbols(")")
}

func (e *Engine) compileSymbol() {
	e.expectOneOfSymbols("(", "-", "~")

	if e.tokenizer.Symbol() == "(" {
		e.advance()

		e.compileExpression()
		e.expectOneOfSymbols(")")

		e.advance()

		return
	}

	operation := e.tokenizer.Symbol()

	e.advance()
	e.compileTerm()

	e.vm.WriteUnaryOperation(operation)

	return
}

// compileExpressionList compiles a (possibly empty) comma-separated list of expressions.
func (e *Engine) compileExpressionList() (expressions int) {
	if !(e.isCurrentSymbol(")")) {
		e.compileExpression()
		expressions++

		for e.isCurrentSymbol(",") {
			e.advance()
			e.compileExpression()
			expressions++
		}
	}

	return expressions
}
