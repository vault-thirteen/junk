package models

import "strings"

// Argument is an argument of a function.
type Argument struct {
	Name Parameter
	Type Parameter
}

// CText returns argument's text as in C language.
func (a Argument) CText() string {
	if len(a.Name.Original) == 0 {
		// Only type is set.
		return a.Type.Original
	}

	// Both type and name are set.
	return a.Type.Original + " " + a.Name.Original
}

// GoText returns argument's text as in Go language.
func (a Argument) GoText() string {
	if len(a.Name.Original) == 0 {
		// Only type is set.
		return a.Type.Original
	}

	// Both type and name are set.
	return a.Name.Original + " " + a.Type.Original
}

// HasPointerType tells whether the argument has a pointer type or a value
// type.
func (a Argument) HasPointerType() bool {
	return strings.Contains(a.Type.Original, "*")
}
