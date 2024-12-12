@ECHO OFF

SET BuildFolder=build
SET App=Manager.exe

SET DbFolder=..\Db
SET DbFile=db.sqlite

%BuildFolder%\%App% -mode=create -database=%DbFolder%\%DbFile%
