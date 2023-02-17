package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// https://docs.min.io/docs/how-to-use-aws-sdk-for-go-with-minio-server.html.
const (
	S3TokenForMinio                 = ""
	S3RegionForMinio                = "us-east-1"
	S3DisableSslForMinio            = true
	S3ForcePathStyleForMinio        = true
	StorageFileDownloadTimeout      = time.Minute * 5
	StorageFileSizeTimeout          = time.Second * 30
	StorageFilePartUploadTimeout    = time.Minute * 1
	StorageFileExistsTimeout        = time.Minute * 1
	StorageListBucketsTimeout       = time.Minute * 5
	StorageReadinessWaitIntervalSec = 5
)

type Storage struct {
	S3ServerAddress    string `split_words:"true" default:"localhost"`
	S3AccessKey        string `split_words:"true"`
	S3Secret           string `split_words:"true"`
	S3Token            string `split_words:"true" default:""`
	S3Region           string `split_words:"true"`
	S3DisableSsl       bool   `split_words:"true" default:"false"`
	S3ForcePathStyle   bool   `split_words:"true" default:"false"`
	S3IsMinio          bool   `split_words:"true" default:"false"`
	S3LocalFilesFolder string `split_words:"true" default:"."`
}

func NewStorage(envPrefix string) (cfg *Storage, err error) {
	cfg = new(Storage)
	err = envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (s *Storage) IsValid() (bool, error) {
	if len(s.S3ServerAddress) < 1 {
		return false, ErrAddress
	}

	if len(s.S3AccessKey) < 1 {
		return false, ErrAccessKey
	}

	if len(s.S3Secret) < 1 {
		return false, ErrSecret
	}

	if !s.S3IsMinio {
		if len(s.S3Region) < 1 {
			return false, ErrRegion
		}
	}

	if len(s.S3LocalFilesFolder) < 1 {
		return false, ErrLocalFilesFolder
	}

	return true, nil
}

func GetStorageSettings() (storageSettings *Storage, err error) {
	storageSettings, err = NewStorage(EnvironmentVariablePrefixApplication)
	if err != nil {
		return nil, err
	}

	_, err = storageSettings.IsValid()
	if err != nil {
		return nil, err
	}

	if storageSettings.S3IsMinio {
		storageSettings.S3Token = S3TokenForMinio
		storageSettings.S3Region = S3RegionForMinio
		storageSettings.S3DisableSsl = S3DisableSslForMinio
		storageSettings.S3ForcePathStyle = S3ForcePathStyleForMinio

		_, err = storageSettings.IsValid()
		if err != nil {
			return nil, err
		}
	}

	return storageSettings, nil
}
