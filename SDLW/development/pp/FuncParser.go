package pp

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/vault-thirteen/SDLW/development/pp/models"
	ae "github.com/vault-thirteen/auxie/errors"
	"github.com/vault-thirteen/auxie/reader"
)

// FuncParser is a C function parser.
type FuncParser struct {
	ignoredWords []models.Word
}

// NewFuncParser is a FuncParser's constructor.
func NewFuncParser(fps *models.FuncProcessorSettings) (fp *FuncParser, err error) {
	if fps == nil {
		return nil, errors.New(models.ErrNoSettings)
	}

	fp = &FuncParser{
		ignoredWords: fps.IgnoredWords,
	}

	return fp, nil
}

// ParseCFunctions reads a file containing definitions of C functions and
// parses its lines into function data. Lines which do not have a simple
// function header are ignored. Function definition must be on a single line.
//
// N.B.
// Golang currently does not support callbacks and C function pointers, so
// implementation of such things in this pp would be a waste of time.
func (fp *FuncParser) ParseCFunctions(filePath string) (funcData []*models.FuncData, err error) {
	// File.
	var f *os.File
	f, err = os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		derr := f.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	// Result Holder.
	funcData = make([]*models.FuncData, 0)
	var fd *models.FuncData

	// Line Reader.
	lr := reader.New(f)
	var line []byte
	line, err = lr.ReadLineEndingWithCRLF()
	for {
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			} else {
				return nil, err
			}
		}

		fd, err = fp.parseCLine(strings.TrimSpace(string(line)))
		// Clear error for special cases.
		if err != nil {
			if strings.HasPrefix(err.Error(), models.ErrNoHeader) {
				err = nil
				log.Print(fmt.Sprintf("ignoring line without header: %s", string(line)))
			} else if strings.HasPrefix(err.Error(), models.ErrNoArgs) {
				err = nil
				log.Print(fmt.Sprintf("ignoring line without arguments: %s", string(line)))
			}
		}
		if err != nil {
			return nil, err
		}

		if fd != nil {
			funcData = append(funcData, fd)
		}

		// Next Line.
		line, err = lr.ReadLineEndingWithCRLF()
	}

	return funcData, nil
}

// parseCLine parses C function data from a text line.
// On error, a null pointer is returned.
// On success, a set pointer is returned.
func (fp *FuncParser) parseCLine(line string) (fd *models.FuncData, err error) {
	fd = &models.FuncData{
		C: &models.FuncDatum{
			SourceCode: []string{line},
		},
		Go: &models.FuncDatum{},
	}

	// Remove words which are ignored, detach stars, remove double spaces.
	clearLine := models.StringWithoutWords(line, fp.ignoredWords)
	clearLine = DetachPointerMarks(clearLine)
	clearLine = RemoveDoubleSpaces(strings.TrimSpace(clearLine))

	err = fp.parseCHead(clearLine, fd)
	if err != nil {
		return nil, err
	}

	err = fp.parseCArgs(clearLine, fd)
	if err != nil {
		return nil, err
	}

	return fd, nil
}

// parseCHead parses C function head data and saves it into dst.
// Function header consists of an output argument and function name.
func (fp *FuncParser) parseCHead(line string, dst *models.FuncData) (err error) {
	parts := strings.Split(line, "(")
	if len(parts) != 2 {
		return fmt.Errorf(models.ErrFmtNoHeader, line)
	}
	header := parts[0]

	headerParts := strings.Split(header, " ")
	if len(headerParts) < 2 {
		return fmt.Errorf(models.ErrBadHeader, header)
	}
	funcName := headerParts[len(headerParts)-1]
	funcRetTypeParts := headerParts[:len(headerParts)-1]
	funcRetType := GlueParts(funcRetTypeParts)

	dst.C.Name.Original = funcName
	dst.C.ReturnedValues = []*models.Argument{
		{
			Type: models.Parameter{
				Original: funcRetType,
			},
		},
	}

	return nil
}

// parseCArgs parses C function's input arguments and saves them into dst.
func (fp *FuncParser) parseCArgs(line string, dst *models.FuncData) (err error) {
	dst.C.Arguments = make([]*models.Argument, 0)

	parts := strings.Split(line, "(")
	if len(parts) != 2 {
		dst.C.Arguments = nil
		return fmt.Errorf(models.ErrFmtNoHeader, line)
	}
	nonHead := parts[1]

	nonHeadParts := strings.Split(nonHead, ")")
	if len(nonHeadParts) < 1 {
		dst.C.Arguments = nil
		return fmt.Errorf(models.ErrFmtNoArgs, line)
	}
	argsPart := nonHeadParts[0]

	argSections := strings.Split(argsPart, ",")

	var arg *models.Argument
	for _, argSection := range argSections {
		arg, err = parseCArgSection(argSection)
		if err != nil {
			dst.C.Arguments = nil
			return err
		}

		dst.C.Arguments = append(dst.C.Arguments, arg)
	}

	return nil
}

// parseCArgSection parses C argument text.
// Examples of argument texts: void, char * s, int size.
func parseCArgSection(s string) (arg *models.Argument, err error) {
	s = strings.TrimSpace(s)
	parts := strings.Split(s, " ")
	partsCount := len(parts)
	switch partsCount {
	case 1:
		// Only type is set. Example: void.
		arg = &models.Argument{
			Type: models.Parameter{
				Original: parts[0],
			},
		}
		return arg, nil

	case 2:
		// Name & type are set. Example: int size.
		arg = &models.Argument{
			Type: models.Parameter{
				Original: parts[0],
			},
			Name: models.Parameter{
				Original: parts[1],
			},
		}
		return arg, nil

	case 3:
		// Name & pointer type ? Example: char * name.
		if (parts[1] == "*") || (parts[1] == "**") {
			arg = &models.Argument{
				Type: models.Parameter{
					Original: GlueParts([]string{parts[0], parts[1]}),
				},
				Name: models.Parameter{
					Original: parts[2],
				},
			}
			return arg, nil
		} else {
			return nil, fmt.Errorf(models.ErrUnsupportedArgumentText, s)
		}

	default:
		return nil, fmt.Errorf(models.ErrUnsupportedArgumentText, s)
	}
}
