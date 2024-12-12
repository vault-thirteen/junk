package models

import "sort"

const (
	PlaceholderTypeM = "M"
	PlaceholderTypeR = "R"
	PlaceholderTypeT = "T"
	PlaceholderTypeU = "U"
)

type Placeholder struct {
	Type string
	Pos  int
}

func NewPlaceholder(typé string, pos int) (ph *Placeholder) {
	if !IsValidPlaceholderType(typé) {
		return nil
	}

	ph = &Placeholder{
		Type: typé,
		Pos:  pos,
	}
	return ph
}

func IsValidPlaceholderType(typé string) bool {
	switch typé {
	case PlaceholderTypeM,
		PlaceholderTypeR,
		PlaceholderTypeT,
		PlaceholderTypeU:
		return true

	default:
		return false
	}
}

func SortPlaceholdersByPosition(phs []Placeholder) {
	sort.Slice(phs, func(i, j int) bool { return phs[i].Pos < phs[j].Pos })
}

func ExcludeNonExistingPlaceholders(phsIn []Placeholder) (phsOut []Placeholder) {
	phsOut = []Placeholder{}

	for _, ph := range phsIn {
		if ph.Pos >= 0 {
			phsOut = append(phsOut, ph)
		}
	}

	return phsOut
}
