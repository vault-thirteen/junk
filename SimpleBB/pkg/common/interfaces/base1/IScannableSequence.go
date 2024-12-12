package base

type IScannableSequence interface {
	IScannable
	Next() bool // = HasNextValue()
}
