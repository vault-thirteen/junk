#include "UTF.h"


// Convert a single/multi-byte string to a UTF-16 string (16-bit).
// We take advantage of the MultiByteToWideChar function that allows to specify
// the charset of the input string. Caller must free the buffer after usage.
LPWSTR CHARtoWCHAR(LPSTR str, UINT codePage) {
    size_t len = strlen(str) + 1;
    int size_needed = MultiByteToWideChar(codePage, 0, str, len, NULL, 0);
    LPWSTR wstr = (LPWSTR) LocalAlloc(LPTR, sizeof(WCHAR) * size_needed);
    MultiByteToWideChar(codePage, 0, str, len, wstr, size_needed);
    return wstr; // Do not forget to use: LocalFree(x).
}


// Convert a UTF-16 string (16-bit) to a single/multi-byte string.
// We take advantage of the WideCharToMultiByte function that allows to specify
// the charset of the output string. Caller must free the buffer after usage.
LPSTR WCHARtoCHAR(LPWSTR wstr, UINT codePage) {
    size_t len = wcslen(wstr) + 1;
    int size_needed = WideCharToMultiByte(codePage, 0, wstr, len, NULL, 0, NULL, NULL);
    LPSTR str = (LPSTR) LocalAlloc(LPTR, sizeof(CHAR) * size_needed);
    WideCharToMultiByte(codePage, 0, wstr, len, str, size_needed, NULL, NULL);
    return str; // Do not forget to use: LocalFree(x).
}
