package configuration

import "strings"

const (
	ServerLoggerTypeUnknown = 0
	ServerLoggerTypeBuiltIn = 1
)

const (
	ServerLoggerTypeAliasBuiltIn = "built-in"
)

type ServerLoggerType uint8

func NewServerLoggerType(
	loggerType string,
) ServerLoggerType {
	loggerType = strings.ToLower(loggerType)
	switch loggerType {
	case ServerLoggerTypeAliasBuiltIn:
		return ServerLoggerTypeBuiltIn
	}
	return ServerLoggerTypeUnknown
}

func (slt ServerLoggerType) IsValid() bool {
	switch slt {
	case ServerLoggerTypeUnknown:
		return false
	case ServerLoggerTypeBuiltIn:
		return true
	}
	return false
}
