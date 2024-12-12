package models

import (
	"fmt"
	"strings"
)

const (
	FormatStringPlaceholderM = "{" + PlaceholderTypeM + "}"
	FormatStringPlaceholderR = "{" + PlaceholderTypeR + "}"
	FormatStringPlaceholderT = "{" + PlaceholderTypeT + "}"
	FormatStringPlaceholderU = "{" + PlaceholderTypeU + "}"
)

const (
	ErrF_TooManyPlaceholders = "too many placeholders: %s"
)

// FormatString is a format string.
// This model is used in internal processes.
type FormatString struct {
	s string

	hasM bool
	hasR bool
	hasT bool
	hasU bool

	posM int
	posR int
	posT int
	posU int

	placeholders []Placeholder
	typé         string
}

func NewFormatString(s string) (fs *FormatString, err error) {
	var nM, nR, nT, nU int
	nM, nR, nT, nU, err = findPlaceholders(s)
	if err != nil {
		return nil, err
	}

	fs = &FormatString{
		s:    s,
		hasM: nM > 0,
		hasR: nR > 0,
		hasT: nT > 0,
		hasU: nU > 0,
	}

	err = fs.calculateType()
	if err != nil {
		return nil, err
	}

	return fs, nil
}

func findPlaceholders(s string) (nM, nR, nT, nU int, err error) {
	nM, err = findPlaceholder(s, FormatStringPlaceholderM)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	nR, err = findPlaceholder(s, FormatStringPlaceholderR)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	nT, err = findPlaceholder(s, FormatStringPlaceholderT)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	nU, err = findPlaceholder(s, FormatStringPlaceholderU)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return nM, nR, nT, nU, nil
}

func findPlaceholder(s string, ph string) (n int, err error) {
	n = strings.Count(s, ph)
	if n > 1 {
		return n, fmt.Errorf(ErrF_TooManyPlaceholders, ph)
	}
	return n, nil
}

func (fs *FormatString) calculateType() (err error) {
	fs.posM = strings.Index(fs.s, FormatStringPlaceholderM)
	fs.posR = strings.Index(fs.s, FormatStringPlaceholderR)
	fs.posT = strings.Index(fs.s, FormatStringPlaceholderT)
	fs.posU = strings.Index(fs.s, FormatStringPlaceholderU)

	placeholders := []Placeholder{
		*NewPlaceholder(PlaceholderTypeM, fs.posM),
		*NewPlaceholder(PlaceholderTypeR, fs.posR),
		*NewPlaceholder(PlaceholderTypeT, fs.posT),
		*NewPlaceholder(PlaceholderTypeU, fs.posU),
	}

	SortPlaceholdersByPosition(placeholders)
	fs.placeholders = ExcludeNonExistingPlaceholders(placeholders)

	fs.typé = ""
	for _, ph := range fs.placeholders {
		fs.typé = fs.typé + ph.Type
	}

	return nil
}

func (fs *FormatString) String() string {
	return fs.s
}

func (fs *FormatString) Type() string {
	return fs.typé
}

func (fs *FormatString) HasM() bool {
	return fs.hasM
}

func (fs *FormatString) HasR() bool {
	return fs.hasR
}

func (fs *FormatString) HasT() bool {
	return fs.hasT
}

func (fs *FormatString) HasU() bool {
	return fs.hasU
}

func (fs *FormatString) Placeholders() []Placeholder {
	return fs.placeholders
}
