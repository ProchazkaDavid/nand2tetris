package compilation

import (
	"errors"
	"fmt"
	"io"

	"github.com/ProchazkaDavid/nand2tetris/compiler/symbol"
	"github.com/ProchazkaDavid/nand2tetris/compiler/token"
	"github.com/ProchazkaDavid/nand2tetris/compiler/tokenizer"
	"github.com/ProchazkaDavid/nand2tetris/compiler/vm"
)

var (
	errUnexpectedToken = errors.New("unexpected token")
	errNoTokens        = errors.New("no more tokens available")
	errCantAdvance     = errors.New("can't advance the tokenizer")
)

// Engine compiles the class in the input file into the output file.
type Engine struct {
	tokenizer    *tokenizer.Tokenizer
	symbolTable  *symbol.Table
	vm           *vm.Writer
	className    string
	ifCounter    int
	whileCounter int
}

// NewEngine creates a new compilation engine with the given input and output.
func NewEngine(input io.Reader, output io.StringWriter) *Engine {
	return &Engine{
		tokenizer:   tokenizer.New(input),
		symbolTable: symbol.NewSymbolTable(),
		vm:          vm.NewWriter(output),
	}
}

// compileParameterList compiles a (possibly empty) parameter list.
// Doest not handle the enclosing "()". Returns number of parameters.
func (e *Engine) compileParameterList() (parameters int) {
	for ; !e.isCurrentSymbol(")"); e.advance() {
		if e.isCurrentSymbol(",") {
			e.advance()
		}

		e.expectType()
		varType := e.getVariableType()

		e.advance()
		e.expectIdentifier()
		e.symbolTable.Define(e.tokenizer.Identifier(), varType, symbol.Arg)
		parameters++
	}

	return parameters
}

// compileSubroutineBody compiles a subroutine's body.
func (e *Engine) compileSubroutineBody(function string, functionType token.KeywordType, classVariables int) {
	e.expectOneOfSymbols("{")

	e.advance()

	variables := 0
	for ; e.isCurrentKeyword(token.Var); e.advance() {
		variables += e.compileVariableDeclarations()
	}

	e.vm.WriteFunction(function, variables)

	switch functionType {
	case token.Constructor:
		e.vm.WritePush(vm.Constant, classVariables)
		e.vm.WriteCall("Memory.alloc", 1)
		e.vm.WritePop(vm.Pointer, 0)

	case token.Method:
		e.vm.WritePush(vm.Arg, 0)
		e.vm.WritePop(vm.Pointer, 0)
	}

	e.compileStatements()

	e.expectOneOfSymbols("}")

	e.advance()
}

// compileVariableDeclarations compiles a variable declarations.
// Returns number of variables.
func (e *Engine) compileVariableDeclarations() (variables int) {
	e.expectOneOfKeywords(token.Var)

	e.advance()
	e.expectType()

	variableType := e.getVariableType()

	e.advance()
	e.expectIdentifier()

	e.symbolTable.Define(e.tokenizer.Identifier(), variableType, symbol.Var)
	variables++

	e.advance()

	for ; e.isCurrentSymbol(","); e.advance() {
		e.advance()
		e.expectIdentifier()

		e.symbolTable.Define(e.tokenizer.Identifier(), variableType, symbol.Var)
		variables++
	}

	e.expectOneOfSymbols(";")

	return variables
}

// handleError handles errors during compilation
func (e *Engine) handleError(err error) {
	panic(fmt.Errorf("engine: %w", err))
}

// advance advances engine's tokenizer
func (e *Engine) advance() {
	if !e.tokenizer.HasMoreTokens() {
		e.handleError(errNoTokens)
	}

	if e.tokenizer.Advance() != nil {
		e.handleError(errCantAdvance)
	}
}

// getVariableType returns string representation of a variable type.
// Should be called after expectType().
func (e *Engine) getVariableType() string {
	if e.tokenizer.TokenType() == token.Keyword {
		return string(e.tokenizer.Keyword())
	}

	return e.tokenizer.Identifier()
}
