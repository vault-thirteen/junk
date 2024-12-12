# SDLW

A wrapper of _SDL_ library for _Go_ language on _Windows_ O.S.

This library also provides various helper-functions to make usage of the _SDL_ 
library much easier.

**Notes**

* This project is frozen due to bugs in _Go_ language and _SDL_ library.
  * More information can be found in comments to the `Test` function and _C_ 
tests.


* This project uses system calls.
  * _CGO_ does not work properly on _Windows_ O.S.
  * _CGO_ is slower than system calls in _Golang_.

## SDL 2.0 API by Category
https://wiki.libsdl.org/SDL2/APIByCategory  
Note that the official documentation is very poor.

## Examples
Code examples are available in the `example` folder.
