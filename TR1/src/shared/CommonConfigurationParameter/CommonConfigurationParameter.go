package ccp

import (
	"encoding/json"
	"fmt"
)

const (
	ValueType_String             = "string"
	ValueType_StringArray        = "strings"
	ValueType_IntegerNumber      = "integer"
	ValueType_IntegerNumberArray = "integers"
	ValueType_Boolean            = "boolean"
	ValueType_MapStringString    = "map"
)

const (
	Administrator                                     = "administrator"
	AllowNativePasswords                              = "allowNativePasswords"
	BodyTemplate_Reg                                  = "bodyTemplate_Reg"
	BodyTemplate_RegRFA                               = "bodyTemplate_RegRFA"
	BodyTemplate_RegApproved                          = "bodyTemplate_RegApproved"
	BodyTemplate_LogIn                                = "bodyTemplate_LogIn"
	BodyTemplate_PwdChange                            = "bodyTemplate_PwdChange"
	BodyTemplate_EmailChange                          = "bodyTemplate_EmailChange"
	CacheControlMaxAge                                = "cacheControlMaxAge"
	CacheRecordTtl                                    = "cacheRecordTtl"
	CacheSizeLimit                                    = "cacheSizeLimit"
	CertFile                                          = "certFile"
	CheckConnLiveness                                 = "checkConnLiveness"
	ClientIPAddressSource_CustomHeader                = "clientIPAddressSourceCustomHeader"
	DatabaseName                                      = "dbName"
	DatabaseType                                      = "databaseType"
	DeveloperMode_HttpHeader_AccessControlAllowOrigin = "devModeHttpHeaderAccessControlAllowOrigin"
	DriverName                                        = "driverName"
	EmailChangeRequestTtl                             = "emailChangeRequestTtl"
	EnableSelfSignedCertificate                       = "enableSelfSignedCertificate"
	FileCacheItemTtl                                  = "fileCacheItemTtl"
	FileCacheSizeLimit                                = "fileCacheSizeLimit"
	FileCacheVolumeLimit                              = "fileCacheVolumeLimit"
	FilesCountToClean                                 = "filesCountToClean"
	Host                                              = "host"
	ImageHeight                                       = "imageHeight"
	ImageWidth                                        = "imageWidth"
	ImagesFolder                                      = "imagesFolder"
	InitConsoleColours                                = "initConsoleColours"
	IsAdminApprovalRequired                           = "isAdminApprovalRequired"
	IsCacheEnabled                                    = "isCacheEnabled"
	IsDatabaseInitialisationUsed                      = "isDatabaseInitialisationUsed"
	IsDebugMode                                       = "isDebugMode"
	IsDeveloperMode                                   = "isDeveloperMode"
	IsImageCleanupAtStartUsed                         = "isImageCleanupAtStartUsed"
	IsImageServerEnabled                              = "isImageServerEnabled"
	IsImageStorageUsed                                = "isImageStorageUsed"
	IsStorageCleaningEnabled                          = "isStorageCleaningEnabled"
	KeyFile                                           = "keyFile"
	LogInRequestTtl                                   = "logInRequestTtl"
	LogOutRequestTtl                                  = "logOutRequestTtl"
	MaxAllowedPacket                                  = "maxAllowedPacket"
	MessageEditTime                                   = "messageEditTime"
	Moderator                                         = "moderator"
	Name                                              = "name"
	Net                                               = "net"
	PageSize                                          = "pageSize"
	Params                                            = "params"
	Password                                          = "password"
	PasswordChangeRequestTtl                          = "passwordChangeRequestTtl"
	Path                                              = "path"
	Port                                              = "port"
	PrivateKeyFilePath                                = "privateKeyFilePath"
	PublicKeyFilePath                                 = "publicKeyFilePath"
	PublicSettingsVersion                             = "publicSettingsVersion"
	PublicSettingsTtl                                 = "publicSettingsTtl"
	RecordCacheItemTtl                                = "recordCacheItemTtl"
	RecordCacheSizeLimit                              = "recordCacheSizeLimit"
	RegistrationRequestTtl                            = "registrationRequestTtl"
	RequestIdLength                                   = "requestIdLength"
	RootFolderPath                                    = "rootFolderPath"
	Schema                                            = "schema"
	SessionMaxDuration                                = "sessionMaxDuration"
	SigningMethod                                     = "signingMethod"
	SiteDomain                                        = "siteDomain"
	SiteName                                          = "siteName"
	SubjectTemplate_VC                                = "subjectTemplate_VC"
	SubjectTemplate_Reg                               = "subjectTemplate_Reg"
	User                                              = "user"
	UserAgent                                         = "userAgent"
	UserNameMaxLenInBytes                             = "userNameMaxLenInBytes"
	UserPasswordMaxLenInBytes                         = "userPasswordMaxLenInBytes"
	VerificationCodeLength                            = "verificationCodeLength"
)

const (
	ErrF_UnknownConfigurationParameterType = "unknown configuration parameter type: %s"
)

type CommonConfigurationParameter struct {
	Name  string
	Type  string
	Value any
}

func ParseCommonConfigurationParameterValue(rt string, rv json.RawMessage) (v any, err error) {
	switch rt {
	case ValueType_String:
		return parseVariableData[string](rv)

	case ValueType_StringArray:
		return parseVariableData[[]string](rv)

	case ValueType_IntegerNumber:
		return parseVariableData[int](rv)

	case ValueType_IntegerNumberArray:
		return parseVariableData[[]int](rv)

	case ValueType_Boolean:
		return parseVariableData[bool](rv)

	case ValueType_MapStringString:
		return parseVariableData[map[string]string](rv)

	default:
		return nil, fmt.Errorf(ErrF_UnknownConfigurationParameterType, rt)
	}
}

func parseVariableData[T any](src json.RawMessage) (dst T, err error) {
	err = json.Unmarshal(src, &dst)
	if err != nil {
		return dst, err
	}
	return dst, nil
}
