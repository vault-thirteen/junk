@ECHO OFF

SET BuildFolder=build
SET SrcFolder=src

rmdir /s %BuildFolder%
cd %SrcFolder%
mkdir ..\%BuildFolder%
go build -o ..\%BuildFolder%\Manager.exe
cd ..
