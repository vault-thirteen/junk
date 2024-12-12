@ECHO OFF

SET BuildFolder=build
SET App=Manager.exe

SET DbFolder=..\Db
SET DbFile=db.sqlite
SET	VideosList=http://example.org:80/data/videos.txt

%BuildFolder%\%App% -mode=import -database=%DbFolder%\%DbFile% -videos=%VideosList%
