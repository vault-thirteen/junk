:: This script guides the installation of the golang driver for SQLite 3.
:: The driver is very old and uses CGO, so its installation is very problematic.

@ECHO OFF

:: 1. Install the 64-bit GCC for Windows O.S. using the TDM-GCC installer.
:: URL: https://jmeubank.github.io/tdm-gcc/download/

:: 2. Install the package.
SET CGO_ENABLED=1
go install github.com/mattn/go-sqlite3@latest
