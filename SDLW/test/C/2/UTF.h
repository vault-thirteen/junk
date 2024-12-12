#include <windows.h>

#ifndef SDL_TEST_UTF_H
#define SDL_TEST_UTF_H

// Convert a UTF-16 string (16-bit) to an OEM string (8-bit)
#define UNICODEtoOEM(str)   WCHARtoCHAR(str, CP_OEMCP)

// Convert an OEM string (8-bit) to a UTF-16 string (16-bit)
#define OEMtoUNICODE(str)   CHARtoWCHAR(str, CP_OEMCP)

// Convert an ANSI string (8-bit) to a UTF-16 string (16-bit)
#define ANSItoUNICODE(str)  CHARtoWCHAR(str, CP_ACP)

// Convert a UTF-16 string (16-bit) to an ANSI string (8-bit)
#define UNICODEtoANSI(str)  WCHARtoCHAR(str, CP_ACP)

LPWSTR CHARtoWCHAR(LPSTR str, UINT codePage);

LPSTR WCHARtoCHAR(LPWSTR wstr, UINT codePage);

#endif //SDL_TEST_UTF_H
