// ElectionFeedback.go.

package eom

import (
	"fmt"
	"strings"
	"time"
)

const (
	ErrfSyntax = "Syntax Error in Record: %v"
)

const ResultTimeFormat = "2006-01-02 15:04:05-07"

type ElectionFeedback struct {
	RawData string
	Result  bool
	Time    time.Time
}

// ParseRawData Method parses the raw Data into more usable internal Values.
func (ef *ElectionFeedback) ParseRawData() error {

	var err error
	var rawDataS1 string
	var parts []string

	// Get Record's Fields.
	rawDataS1 = strings.TrimPrefix(ef.RawData, "(")
	rawDataS1 = strings.TrimSuffix(rawDataS1, ")")
	parts = strings.Split(rawDataS1, ",")
	if len(parts) != 2 {
		return fmt.Errorf(ErrfSyntax, ef.RawData)
	}
	parts[1] = strings.TrimPrefix(parts[1], `"`)
	parts[1] = strings.TrimSuffix(parts[1], `"`)

	// Raw Result -> Boolean.
	switch parts[0] {

	case "t":
		ef.Result = true

	case "f":
		ef.Result = false

	case "true":
		ef.Result = true

	case "false":
		ef.Result = false

	default:
		return fmt.Errorf(ErrfSyntax, ef.RawData)
	}

	// Raw Time -> Time.
	if len(parts[1]) == 0 {
		ef.Time = time.Time{}
	} else {
		ef.Time, err = time.Parse(ResultTimeFormat, parts[1])
		if err != nil {
			return err
		}
	}

	return nil
}
