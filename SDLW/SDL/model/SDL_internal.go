package m

// SDL_VARIABLE_LENGTH_ARRAY
/*
// This is for a variable-length array at the end of a struct:
// struct x { int y; char z[SDL_VARIABLE_LENGTH_ARRAY]; };
// Use this because GCC 2 needs different magic than other compilers.
#if (defined(__GNUC__) && (__GNUC__ <= 2)) || defined(__CC_ARM) || defined(__cplusplus)
#define SDL_VARIABLE_LENGTH_ARRAY 1
#else
#define SDL_VARIABLE_LENGTH_ARRAY
#endif
*/
const SDL_VARIABLE_LENGTH_ARRAY = 1 //TODO: ???

type Void = byte // Void type is used only as a pointer, i.e. *void.
type Char = byte

type Int = int32
type Uint = uint32

type Int8 = int8
type Uint8 = uint8

type Int16 = int16
type Uint16 = uint16

type Int32 = int32
type Uint32 = uint32

type Float = float32
type Double = float64

type SizeT = int64
