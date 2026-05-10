::============================================================================::
:: This script must be started from its folder ::
::============================================================================::
@ECHO OFF

SET BUILD_DIR=_BUILD_

FOR %%f IN ("%CD%") DO SET LastPathElement=%%~nxf
SET CUR_FOLDER_NAME=%LastPathElement%
IF "%CUR_FOLDER_NAME%" == "script" ( ECHO Welcome ) ELSE (
    ECHO This script must be started from its folder. Press any key to exit.
    EXIT /B 1
)

:: CD to root folder.
CD ..
MKDIR %BUILD_DIR%
MKDIR %BUILD_DIR%\service

CALL :BuildServiceExecutable "AuthService"
IF %ERRORLEVEL% NEQ 0 ( GOTO :BadExit )

CALL :BuildServiceExecutable "CaptchaService"
IF %ERRORLEVEL% NEQ 0 ( GOTO :BadExit )

CALL :BuildServiceExecutable "GatewayService"
IF %ERRORLEVEL% NEQ 0 ( GOTO :BadExit )

CALL :BuildServiceExecutable "MailerService"
IF %ERRORLEVEL% NEQ 0 ( GOTO :BadExit )

CALL :BuildServiceExecutable "MessageService"
IF %ERRORLEVEL% NEQ 0 ( GOTO :BadExit )

MKDIR %BUILD_DIR%\tool

CALL :BuildToolExecutable "MakeJWToken"
IF %ERRORLEVEL% NEQ 0 ( GOTO :BadExit )

ECHO Copying files ...
XCOPY assets %BUILD_DIR%\assets /S/I/Q
XCOPY cert %BUILD_DIR%\cert /S/I/Q
XCOPY config %BUILD_DIR%\config /S/I/Q
XCOPY script\start_service_*.bat %BUILD_DIR%\ /Q
XCOPY script\start_tool_*.bat %BUILD_DIR%\ /Q

MKDIR %BUILD_DIR%\captcha

EXIT /B 0

::============================================================================::

:BuildServiceExecutable
SET SERVICE_NAME=%~1
ECHO Building service %SERVICE_NAME%
CD src\services\%SERVICE_NAME%\
go build -o ..\..\..\%BUILD_DIR%\service\%SERVICE_NAME%\
IF %ERRORLEVEL% NEQ 0 EXIT /B %ERRORLEVEL%
CD ..\..\..\
EXIT /B 0

:BuildToolExecutable
SET TOOL_NAME=%~1
ECHO Building tool %TOOL_NAME%
CD src\tool\%TOOL_NAME%\
go build -o ..\..\..\%BUILD_DIR%\tool\
IF %ERRORLEVEL% NEQ 0 EXIT /B %ERRORLEVEL%
CD ..\..\..\
EXIT /B 0

:BadExit
PAUSE
EXIT /B %ERRORLEVEL%

::============================================================================::
