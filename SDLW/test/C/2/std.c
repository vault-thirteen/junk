#include "std.h"


// Creates a string and fills it with a string formatted with the specified
// format and arguments. Caller must free the string after usage.
char *smprintf(const char *fmt, va_list args) {
    // 1. Measure the required size of buffer.
    va_list args2;
    va_copy(args2, args);
    int len = vsnprintf(NULL, 0, fmt, args2);
            va_end(args2);
    len++;

    // 2. Create the buffer.
    char *buf = malloc(len);
    if (buf == NULL) {
        return NULL;
    }

    // 3. Fill the buffer.
    vsnprintf(buf, len, fmt, args);
    return buf; // Do not forget to use: free(x).
}
