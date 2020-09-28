package compilation

import (
	"fmt"

	"github.com/ProchazkaDavid/nand2tetris/compiler/token"
	"github.com/ProchazkaDavid/nand2tetris/compiler/vm"
)

// compileStatements compiles a sequence of statements.
// Does not handle the enclosing "{}".
func (e *Engine) compileStatements() {
	for e.isOneOfKeywords(token.Let, token.If, token.While, token.Do, token.Return) {
		switch e.tokenizer.Keyword() {
		case token.Let:
			e.compileLet()
		case token.If:
			e.compileIf()
		case token.While:
			e.compileWhile()
		case token.Do:
			e.compileDo()
		case token.Return:
			e.compileReturn()
		}
	}
}

// compileLet compiles a let statement.
func (e *Engine) compileLet() {
	e.expectOneOfKeywords(token.Let)

	e.advance()
	e.expectIdentifier()

	variableName := e.tokenizer.Identifier()

	isArray := false

	e.advance()
	if e.isCurrentSymbol("[") {
		isArray = true

		e.advance()
		e.compileExpression()

		e.expectOneOfSymbols("]")

		e.vm.WritePush(vm.GetSegment(e.symbolTable.KindOf(variableName)), e.symbolTable.IndexOf(variableName))
		e.vm.WriteArithmetic("+")

		e.advance()
	}

	e.expectOneOfSymbols("=")
	e.advance()

	e.compileExpression()

	if isArray {
		e.vm.WritePop(vm.Temp, 0)
		e.vm.WritePop(vm.Pointer, 1)
		e.vm.WritePush(vm.Temp, 0)
		e.vm.WritePop(vm.That, 0)
	} else {
		e.vm.WritePop(vm.GetSegment(e.symbolTable.KindOf(variableName)), e.symbolTable.IndexOf(variableName))
	}

	e.expectOneOfSymbols(";")

	e.advance()
}

// compileIf compiles a if statement, possibly with a trailing else clause.
func (e *Engine) compileIf() {
	trueLabel := fmt.Sprintf("IF_TRUE%d", e.ifCounter)
	falseLabel := fmt.Sprintf("IF_FALSE%d", e.ifCounter)
	endLabel := fmt.Sprintf("IF_END%d", e.ifCounter)

	e.ifCounter++

	e.expectOneOfKeywords(token.If)

	e.advance()
	e.expectOneOfSymbols("(")

	e.advance()
	e.compileExpression()

	e.expectOneOfSymbols(")")

	e.vm.WriteIf(trueLabel)
	e.vm.WriteGoto(falseLabel)
	e.vm.WriteLabel(trueLabel)

	e.advance()
	e.expectOneOfSymbols("{")

	e.advance()
	e.compileStatements()

	e.expectOneOfSymbols("}")

	e.advance()

	if e.isCurrentKeyword(token.Else) {
		e.vm.WriteGoto(endLabel)
	}

	e.vm.WriteLabel(falseLabel)

	if e.isCurrentKeyword(token.Else) {
		e.expectOneOfKeywords(token.Else)
		e.advance()
		e.expectOneOfSymbols("{")

		e.advance()
		e.compileStatements()

		e.expectOneOfSymbols("}")

		e.advance()
		e.vm.WriteLabel(endLabel)
	}
}

// compileWhile compiles a while statement.
func (e *Engine) compileWhile() {
	expressionLabel := fmt.Sprintf("WHILE_EXP%d", e.whileCounter)
	endLabel := fmt.Sprintf("WHILE_END%d", e.whileCounter)

	e.whileCounter++

	e.expectOneOfKeywords(token.While)

	e.vm.WriteLabel(expressionLabel)

	e.advance()
	e.expectOneOfSymbols("(")

	e.advance()
	e.compileExpression()

	e.expectOneOfSymbols(")")

	e.vm.WriteArithmetic("~")
	e.vm.WriteIf(endLabel)

	e.advance()
	e.expectOneOfSymbols("{")

	e.advance()
	e.compileStatements()

	e.expectOneOfSymbols("}")

	e.advance()

	e.vm.WriteGoto(expressionLabel)
	e.vm.WriteLabel(endLabel)
}

// compileDo compiles a do statement.
func (e *Engine) compileDo() {
	e.expectOneOfKeywords(token.Do)

	e.advance()
	e.expectIdentifier()

	identifier := e.tokenizer.Identifier()

	e.advance()

	expressions := 0
	isCurrentClassCall := false
	class := ""
	method := ""

	if e.isCurrentSymbol("(") {
		isCurrentClassCall = true

		class = e.className
		method = identifier

		expressions++
	} else {
		class = identifier

		e.expectOneOfSymbols(".")

		e.advance()
		e.expectIdentifier()

		method = e.tokenizer.Identifier()

		if classType, ok := e.symbolTable.TypeOf(identifier); ok {
			e.vm.WritePush(vm.GetSegment(e.symbolTable.KindOf(identifier)), e.symbolTable.IndexOf(identifier))
			expressions++

			class = classType
			method = e.tokenizer.Identifier()
		}

		e.advance()
	}

	e.expectOneOfSymbols("(")

	e.advance()

	if isCurrentClassCall {
		e.vm.WritePush(vm.Pointer, 0)
	}

	expressions += e.compileExpressionList()

	e.expectOneOfSymbols(")")
	e.advance()

	e.expectOneOfSymbols(";")
	e.advance()

	e.vm.WriteCall(fmt.Sprintf("%s.%s", class, method), expressions)
	e.vm.WritePop(vm.Temp, 0)
}

// compileReturn compiles a return statement.
func (e *Engine) compileReturn() {
	e.expectOneOfKeywords(token.Return)

	e.advance()
	isNakedReturn := e.isCurrentSymbol(";")

	if !isNakedReturn {
		e.compileExpression()
	}

	e.expectOneOfSymbols(";")
	e.advance()

	if isNakedReturn {
		e.vm.WritePush(vm.Constant, 0)
	}

	e.vm.WriteReturn()
}
