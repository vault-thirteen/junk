package m

// Bool
// SDL_stdinc.h
type Bool Int

const (
	FALSE = Bool(0)
	TRUE  = Bool(1)
)

func BoolFromUintptr(x uintptr) bool {
	if int(x) == 0 {
		return false
	}
	if int(x) == 1 {
		return true
	}
	panic(x)
}

func BoolToUintptr(b bool) uintptr {
	if b {
		return uintptr(1)
	}
	return uintptr(0)
}
