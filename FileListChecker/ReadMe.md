## File List Checker  

A tool to check existence of files listed in a file.

A small tool which reads a text file containing a list of file names and 
ensures that all the listed files do exist. During the check a path prefix is 
appended to each of the file names in the list.

### Note

Please, note that functionality of this tool is now implemented in the `Hasher` 
tool which is available at the following repository: 
https://github.com/vault-thirteen/Hasher

### Command line arguments
    1-st:	Folder used as a prefix for each file in the list.
    2-nd:	Path to a file containing the list of file names.

### Usage Example
> FileListChecker.exe test test/Names.txt

You should see something similar to:  
>    File does not exist: test\ccc.txt

### Installation
`go install github.com/vault-thirteen/FileListChecker/cmd/checker@latest`
