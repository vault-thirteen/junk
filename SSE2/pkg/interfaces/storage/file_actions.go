package storage

import (
	"context"

	"github.com/vault-thirteen/junk/SSE2/pkg/models/storage"
)

type FileActions interface {
	GetS3LocalFilesFolder() (folderPath string)

	GetFileSize(
		ctx context.Context,
		srcBucket string,
		srcPath string,
	) (fileSize int, err error)

	DownloadFile(
		ctx context.Context,
		srcBucket string,
		srcPath string,
		dstLocalFolderPath string,
	) (result *storage.DownloadResult, err error)

	UploadFile(
		ctx context.Context,
		srcLocalFilePath string,
		contentType string,
		dstBucket string,
		dstFilePath string,
	) (err error)

	DoesFileExist(
		ctx context.Context,
		bucket string,
		filePath string,
	) (fileExists *bool, err error)
}
