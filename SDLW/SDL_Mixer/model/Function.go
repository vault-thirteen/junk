package mm

import m "github.com/vault-thirteen/SDLW/SDL/model"

// Unfortunately, developers of SDL were drunk when they called function
// pointer 'function', so, we can not say for sure what this is.
// int (SDLCALL *function)(const char*, void*)
type Function *func(*m.Char, *m.Void) m.Int
