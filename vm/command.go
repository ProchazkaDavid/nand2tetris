package main

// commandType represents command type
type commandType int

const (
	// arithmeticCmd command
	arithmeticCmd commandType = iota
	// pushCmd command
	pushCmd
	// popCmd command
	popCmd
	// labelCmd command
	labelCmd
	// gotoCmd command
	gotoCmd
	// ifCmd command
	ifCmd
	// functionCmd command
	functionCmd
	// returnCmd command
	returnCmd
	// callCmd command
	callCmd
)
