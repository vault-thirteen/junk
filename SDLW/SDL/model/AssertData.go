package m

// AssertData
/*
typedef struct SDL_AssertData
{
	int always_ignore;
	unsigned int trigger_count;
	const char *condition;
	const char *filename;
	int linenum;
	const char *function;
	const struct SDL_AssertData *next;
} SDL_AssertData;
*/
// SDL_assert.h
type AssertData struct {
	AlwaysIgnore Int
	TriggerCount Uint
	Condition    *Char
	Filename     *Char
	LineNum      Int
	Function     *Char
	Next         *AssertData
}
