// ElectionType.go.

package eom

const ElectionTypeSingleMaster = 1

type ElectionType int

func (et ElectionType) IsValid() bool {

	switch et {

	case ElectionTypeSingleMaster:
		return true

	default:
		return false
	}
}
