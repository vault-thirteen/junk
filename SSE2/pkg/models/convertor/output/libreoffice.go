package output

import (
	"strings"

	"github.com/pkg/errors"
)

const (
	ActionPrefix             = "convert"
	SeparatorArrow           = "->"
	FilterMention            = "using filter"
	SeparatorColon           = ":"
	ArrowSeparatedPartsCount = 2
	ColonSeparatedPartsCount = 2
)

var (
	ErrNoLinesToProcess          = errors.New("no lines to process")
	ErrFArrowSeparatedPartsCount = "arrow separated parts count mismatch, %v vs %v"
	ErrFColonSeparatedPartsCount = "colon separated parts count mismatch, %v vs %v"
)

func ParseLibreOfficeConverterOutput(outputLines []string) (output *Output, err error) {
	if len(outputLines) < 1 {
		return nil, ErrNoLinesToProcess
	}

	arrowSeparatedParts := strings.Split(outputLines[0], SeparatorArrow)
	if len(arrowSeparatedParts) != ArrowSeparatedPartsCount {
		return nil, errors.Errorf(ErrFArrowSeparatedPartsCount, ArrowSeparatedPartsCount, len(arrowSeparatedParts))
	}

	output = new(Output)

	output.SourceFilePath = strings.TrimSpace(
		strings.TrimPrefix(arrowSeparatedParts[0], ActionPrefix),
	)

	tmpParts := strings.Split(arrowSeparatedParts[1], SeparatorColon)
	if len(tmpParts) != ColonSeparatedPartsCount {
		return nil, errors.Errorf(ErrFColonSeparatedPartsCount, ColonSeparatedPartsCount, len(tmpParts))
	}

	output.MiscellaneousData = strings.TrimSpace(tmpParts[1])

	output.DestinationFilePath = strings.TrimSpace(
		strings.TrimSuffix(
			strings.TrimSpace(tmpParts[0]),
			FilterMention,
		),
	)

	return output, nil
}
