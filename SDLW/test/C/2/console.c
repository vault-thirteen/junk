#include "console.h"

// Prints the formatted text in console.
void console_print(const char *fmt, ...) {
    va_list args;
            va_start(args, fmt);
    char *buf = smprintf(fmt, args);
            va_end(args);
    if (buf == NULL) {
        return;
    }

    PDWORD cChars = NULL;
    HANDLE std_out = GetStdHandle(STD_OUTPUT_HANDLE);
    LPWSTR ws = CHARtoWCHAR(buf, CP_ACP);
    WriteConsoleW(std_out, ws, wcslen(ws), cChars, NULL);

    LocalFree(ws);
    free(buf);
}
