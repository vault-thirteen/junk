package models

const (
	ErrNoSettings              = "no settings"
	ErrNoFuncData              = "no function data"
	ErrNoTypeMapping           = "no type mapping: %s"
	ErrNoHeader                = "no header"
	ErrFmtNoHeader             = ErrNoHeader + ": %s"
	ErrBadHeader               = "bad header: %s"
	ErrNoArgs                  = "no args"
	ErrFmtNoArgs               = ErrNoArgs + ": %s"
	ErrUnsupportedArgumentText = "unsupported argument text: %s"
	ErrMultipleReturnedValues  = "multiple returned values"
	ErrMappingSyntaxError      = "mapping syntax error"
	ErrDuplicateKey            = "duplicate key: %s"
)
