package command

// Type represents command type
type Type int

const (
	// Arithmetic command
	Arithmetic Type = iota
	// Push command
	Push
	// Pop command
	Pop
	// Label command
	Label
	// Goto command
	Goto
	// If command
	If
	// Function command
	Function
	// Return command
	Return
	// Call command
	Call
)
