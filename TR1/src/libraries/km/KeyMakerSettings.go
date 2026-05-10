package km

type KeyMakerSettings struct {
	SigningMethodName  string
	PrivateKeyFilePath string
	PublicKeyFilePath  string
	IsCacheEnabled     bool
	CacheSizeLimit     int
	CacheRecordTtl     uint
}
