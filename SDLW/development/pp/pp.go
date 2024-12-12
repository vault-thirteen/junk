package pp

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/vault-thirteen/SDLW/development/pp/models"
)

const (
	CR      = "\r"
	LF      = "\n"
	NewLine = CR + LF
)

// GetFuncDataFromFile gets information about C functions and converts it into
// Go functions' data. This function is an entry point into this package.
func GetFuncDataFromFile(filePath string, fps *models.FuncProcessorSettings) (fds []*models.FuncData, err error) {
	log.Println(fmt.Sprintf("Parsing a file: %s.", filePath))

	var fp *FuncParser
	fp, err = NewFuncParser(fps)
	if err != nil {
		return nil, err
	}

	fds, err = fp.ParseCFunctions(filePath)
	if err != nil {
		return nil, err
	}

	var fproc *FuncProcessor
	fproc, err = NewFuncProcessor(fps)
	if err != nil {
		return nil, err
	}

	err = fproc.CToGo(fds)
	if err != nil {
		return nil, err
	}

	return fds, nil
}

// ListGoFuncNames lists names of Go functions.
func ListGoFuncNames(fds []*models.FuncData) {
	for _, fd := range fds {
		fmt.Println(fd.Go.Name.Original)
	}
}

// PrintGoSysCallWrapperDrafts prints a draft of Go wrapper-functions using
// system calls. This output should be manually checked and corrected where
// needed.
func PrintGoSysCallWrapperDrafts(fds []*models.FuncData) (err error) {
	var buf string
	for _, fd := range fds {
		buf, err = composeGoWrapperFuncDraft(fd)
		if err != nil {
			return err
		}

		fmt.Println(buf)
		fmt.Println()
	}

	return nil
}

// composeGoWrapperFuncDraft composes a draft of Go wrapper functions.
func composeGoWrapperFuncDraft(fd *models.FuncData) (s string, err error) {
	if fd == nil {
		return s, errors.New(models.ErrNoFuncData)
	}

	var sb strings.Builder

	for _, scLine := range fd.C.SourceCode {
		sb.WriteString(`//` + scLine + NewLine)
	}

	if len(fd.Go.ReturnedValues) > 1 {
		return "", errors.New(models.ErrMultipleReturnedValues)
	}

	hasRetVal := len(fd.Go.ReturnedValues) > 0

	// First Line ('func').
	sb.WriteString("func" + " " + fd.Go.Name.Original + "(")
	sb.WriteString(goFuncArgsText(fd.Go.Arguments))
	sb.WriteString(")")
	sb.WriteString("(")
	sb.WriteString(goFuncArgsText(fd.Go.ReturnedValues))
	sb.WriteString(") {" + NewLine)

	// System Call Line.
	if hasRetVal {
		sb.WriteString("ret, _, _ := ")
	} else {
		sb.WriteString("_, _, _ = ")
	}
	sb.WriteString("syscall.SyscallN(" + "fn" + fd.Go.Name.Original + ", ")
	for _, arg := range fd.Go.Arguments { // List of input arguments.
		if arg.Type.Original == TypeNameString {
			sb.WriteString("uintptr(unsafe.Pointer(BytePtrFromStringP(" + arg.Name.Original + "))), ")
		} else {
			if arg.HasPointerType() {
				sb.WriteString("uintptr(unsafe.Pointer(" + arg.Name.Original + ")), ")
			} else {
				sb.WriteString("uintptr(" + arg.Name.Original + "), ")
			}
		}
	}
	sb.WriteString(")" + NewLine)

	// Return Statement Line.
	if hasRetVal {
		if fd.Go.ReturnedValues[0].Type.Original == TypeNameString {
			sb.WriteString("return windows.BytePtrToString((*byte)(unsafe.Pointer(ret)))" + NewLine)
		} else {
			if fd.Go.ReturnedValues[0].HasPointerType() {
				sb.WriteString("return (" + fd.Go.ReturnedValues[0].Type.Original + ")(unsafe.Pointer(ret))" + NewLine)
			} else {
				sb.WriteString("return (" + fd.Go.ReturnedValues[0].Type.Original + ")(ret)" + NewLine)
			}
		}
	}

	sb.WriteString("}")

	return sb.String(), nil
}

// goFuncArgsText returns argument's names and types as in Go language.
func goFuncArgsText(args []*models.Argument) (s string) {
	var setArgs = make([]*models.Argument, 0, len(args))
	for _, arg := range args {
		if arg != nil {
			setArgs = append(setArgs, arg)
		}
	}

	if len(setArgs) == 0 {
		return s
	}

	argTexts := make([]string, 0, len(setArgs))
	for _, arg := range setArgs {
		argTexts = append(argTexts, arg.GoText())
	}

	return strings.Join(argTexts, ", ")
}
