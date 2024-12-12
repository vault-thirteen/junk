package pp

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vault-thirteen/SDLW/development/pp/models"
)

const (
	TypeNameVoid   = "void"
	TypeNameString = "string"
)

// FuncProcessor is a C function processor.
type FuncProcessor struct {
	funcNamePrefix   string
	cToGoTypeMapping map[string]string
}

// NewFuncProcessor is a FuncProcessor's constructor.
func NewFuncProcessor(fps *models.FuncProcessorSettings) (fp *FuncProcessor, err error) {
	if fps == nil {
		return nil, errors.New(models.ErrNoSettings)
	}

	fp = &FuncProcessor{
		funcNamePrefix:   fps.FuncNamePrefix,
		cToGoTypeMapping: fps.CToGoTypeMapping,
	}

	return fp, nil
}

// CToGo fills Go functions using information about C functions.
// C functions must be parsed and ready.
func (fp *FuncProcessor) CToGo(fds []*models.FuncData) (err error) {
	for _, fd := range fds {
		if fd == nil {
			continue
		}
		if fd.C == nil {
			continue
		}

		fd.Go = &models.FuncDatum{
			Name: models.Parameter{
				Original: strings.TrimPrefix(fd.C.Name.Original, fp.funcNamePrefix),
			},
		}

		fd.Go.ReturnedValues, err = fp.cToGoReturnedValues(fd)
		if err != nil {
			return err
		}

		fd.Go.Arguments, err = fp.cToGoArguments(fd)
		if err != nil {
			return err
		}
	}

	return nil
}

// cToGoReturnedValues converts C returned values into Go returned values.
func (fp *FuncProcessor) cToGoReturnedValues(fd *models.FuncData) (goRetValues []*models.Argument, err error) {
	if fd == nil {
		return nil, errors.New(models.ErrNoFuncData)
	}

	goRetValues = make([]*models.Argument, 0, len(fd.C.ReturnedValues))
	var goRV *models.Argument
	var typeMappingExists bool
	for _, cRV := range fd.C.ReturnedValues {
		// Golang does not use the `void` output type.
		if cRV.Type.Original == TypeNameVoid {
			continue
		}

		_, typeMappingExists = fp.cToGoTypeMapping[cRV.Type.Original]
		if !typeMappingExists {
			return nil, fmt.Errorf(models.ErrNoTypeMapping, cRV.Type.Original)
		}

		goRV = &models.Argument{
			Name: models.Parameter{
				Original: cRV.Name.Original,
			},
			Type: models.Parameter{
				Original: fp.cToGoTypeMapping[cRV.Type.Original],
			},
		}

		goRetValues = append(goRetValues, goRV)
	}

	return goRetValues, nil
}

// cToGoArguments converts C arguments into Go arguments.
func (fp *FuncProcessor) cToGoArguments(fd *models.FuncData) (goArgs []*models.Argument, err error) {
	if fd == nil {
		return nil, errors.New(models.ErrNoFuncData)
	}

	goArgs = make([]*models.Argument, 0, len(fd.C.Arguments))
	var goArg *models.Argument
	var typeMappingExists bool
	for _, cArg := range fd.C.Arguments {
		// Golang does not use the `void` input type.
		if cArg.Type.Original == TypeNameVoid {
			continue
		}

		_, typeMappingExists = fp.cToGoTypeMapping[cArg.Type.Original]
		if !typeMappingExists {
			return nil, fmt.Errorf(models.ErrNoTypeMapping, cArg.Type.Original)
		}

		goArg = &models.Argument{
			Name: models.Parameter{
				Original: cArg.Name.Original,
			},
			Type: models.Parameter{
				Original: fp.cToGoTypeMapping[cArg.Type.Original],
			},
		}

		goArgs = append(goArgs, goArg)
	}

	return goArgs, nil
}
