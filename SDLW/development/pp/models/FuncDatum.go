package models

type FuncDatum struct {
	SourceCode     []string
	Name           Parameter
	ReturnedValues []*Argument
	Arguments      []*Argument
}
