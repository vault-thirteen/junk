package configuration

import "strings"

const (
	ServerStorageTypeUnknown = 0
	ServerStorageTypeMysql   = 1
)

const (
	ServerStorageTypeAliasMysql = "mysql"
)

type ServerStorageType uint8

func NewServerStorageType(
	storageType string,
) ServerStorageType {
	storageType = strings.ToLower(storageType)
	switch storageType {
	case ServerStorageTypeAliasMysql:
		return ServerStorageTypeMysql
	}
	return ServerStorageTypeUnknown
}

func (sst ServerStorageType) IsValid() bool {
	switch sst {
	case ServerStorageTypeUnknown:
		return false
	case ServerStorageTypeMysql:
		return true
	}
	return false
}
