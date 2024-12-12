package models

import (
	"errors"
	"fmt"
)

// FuncProcessorSettings is settings of the parser and processor.
type FuncProcessorSettings struct {
	// List of words to be removed from text during parsing of C functions.
	IgnoredWords []Word

	// Function name prefix which will be removed from C function names during
	// the conversion to Go function names.
	FuncNamePrefix string

	// Type mappings.
	// Key = C type name, Value = Go type name.
	CToGoTypeMapping map[string]string
}

// NewFuncProcessorSettings is a constructor of FuncProcessorSettings.
func NewFuncProcessorSettings(ignoredWords []Word, funcNamePrefix string, cToGoTypeMapping map[string]string) (fps *FuncProcessorSettings, err error) {
	fps = &FuncProcessorSettings{
		IgnoredWords:     ignoredWords,
		FuncNamePrefix:   funcNamePrefix,
		CToGoTypeMapping: cToGoTypeMapping,
	}

	err = fps.check()
	if err != nil {
		return nil, err
	}

	return fps, nil
}

// check checks the settings for correctness.
func (fps *FuncProcessorSettings) check() (err error) {
	for _, iw := range fps.IgnoredWords {
		err = iw.Check()
		if err != nil {
			return err
		}
	}

	w := Word(fps.FuncNamePrefix)
	err = w.Check()
	if err != nil {
		return err
	}

	return nil
}

// NewFuncProcessorSettingsFromJson creates FuncProcessorSettings from bytes of
// a JSON file.
func NewFuncProcessorSettingsFromJson(jsonData []byte) (fps *FuncProcessorSettings, err error) {
	var fpsj *FuncProcessorSettingsJson
	fpsj, err = NewFuncProcessorSettingsJson(jsonData)
	if err != nil {
		return nil, err
	}

	var ignoredWords = make([]Word, 0, len(fpsj.IgnoredWords))
	for _, iw := range fpsj.IgnoredWords {
		ignoredWords = append(ignoredWords, Word(iw))
	}

	var funcNamePrefix = fpsj.FuncNamePrefix

	var cToGoTypeMapping = make(map[string]string)
	var isDuplicate bool
	for _, m := range fpsj.CToGoTypeMapping {
		if len(m) != 2 {
			return nil, errors.New(ErrMappingSyntaxError)
		}

		key := m[0]
		value := m[1]

		_, isDuplicate = cToGoTypeMapping[key]
		if isDuplicate {
			return nil, fmt.Errorf(ErrDuplicateKey, key)
		}

		cToGoTypeMapping[key] = value
	}

	fps = &FuncProcessorSettings{
		IgnoredWords:     ignoredWords,
		FuncNamePrefix:   funcNamePrefix,
		CToGoTypeMapping: cToGoTypeMapping,
	}

	return fps, nil
}
