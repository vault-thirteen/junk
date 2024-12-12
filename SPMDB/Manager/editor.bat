@ECHO OFF

SET BuildFolder=build
SET App=Manager.exe

SET DbFolder=..\Db
SET DbFile=db.sqlite

%BuildFolder%\%App% -mode=edit -database=%DbFolder%\%DbFile% -gui=gui
