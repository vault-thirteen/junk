@ECHO OFF

SET BuildFolder=build
SET App=Manager.exe

SET DbFolder=..\Db
SET DbFile=db.sqlite
SET ScriptsFolder1=scripts\sql\init
SET ScriptsFolder2=scripts\sql\data

%BuildFolder%\%App% -mode=initialize -database=%DbFolder%\%DbFile% -scripts=%ScriptsFolder1%
%BuildFolder%\%App% -mode=initialize -database=%DbFolder%\%DbFile% -scripts=%ScriptsFolder2%
