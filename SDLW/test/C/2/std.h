#ifndef SDL_TEST_STD_H
#define SDL_TEST_STD_H

#include <stdarg.h>
#include <stdlib.h>
#include <stdio.h>

// Creates a string and fills it with a string formatted with the specified
// format and arguments. Caller must free the string after usage.
char *smprintf(const char *fmt, va_list args);

#endif //SDL_TEST_STD_H
