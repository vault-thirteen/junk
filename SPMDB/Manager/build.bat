@ECHO OFF

SET BuildFolder=build
SET SrcFolder=src

cd %SrcFolder%
mkdir ..\%BuildFolder%
go build -o ..\%BuildFolder%\Manager.exe
cd ..
