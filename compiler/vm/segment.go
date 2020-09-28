package vm

import "github.com/ProchazkaDavid/nand2tetris/compiler/symbol"

// Segment type
type Segment string

// Possible VM segments and their string representation
const (
	Unknown  Segment = ""
	Constant Segment = "constant"
	Arg      Segment = "argument"
	Local    Segment = "local"
	Static   Segment = "static"
	This     Segment = "this"
	That     Segment = "that"
	Pointer  Segment = "pointer"
	Temp     Segment = "temp"
)

// GetSegment return the corresponding VM segment based on the IdentifierType
// of the given variable
func GetSegment(variableType symbol.Identifier) Segment {
	switch variableType {
	case symbol.Static:
		return Static
	case symbol.Field:
		return This
	case symbol.Arg:
		return Arg
	case symbol.Var:
		return Local
	default:
		return Unknown
	}
}
