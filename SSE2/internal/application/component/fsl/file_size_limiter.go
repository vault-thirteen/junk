package fsl

import (
	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/vault-thirteen/junk/SSE2/internal/fileext"
	"github.com/vault-thirteen/junk/SSE2/pkg/models/file-size-limiter/settings"
	"github.com/vault-thirteen/junk/SSE2/pkg/models/mimetype"
)

const (
	MB                    = 1 * 1000 * 1000
	InternalFileSizeLimit = 1000 * MB // 1 GB.
)

const (
	MsgFDebugConfig = "file size limiter configuration: %+v"

	ErrFMimeTypeNotAvailable            = "mime type '%v' is not available"
	ErrFDuplicateMimeType               = "duplicate mime type '%v'"
	ErrFFileSizeDidNotPassInternalCheck = "file size for mime type '%v' did not pass the internal check"
)

type FileSizeLimiter struct {
	settingsSourceFilePath   string
	fileSizeLimitPerMimeType map[mimetype.Template]FileSize
}

func NewFileSizeLimiter(
	logger *zerolog.Logger,
	settingsSourceFilePath string,
) (limiter *FileSizeLimiter, err error) {
	limiter = new(FileSizeLimiter)

	limiter.settingsSourceFilePath = settingsSourceFilePath

	err = limiter.loadSettings()
	if err != nil {
		return nil, err
	}

	logger.Debug().Msg(pretty.Sprintf(MsgFDebugConfig, limiter.settingsSourceFilePath))

	return limiter, nil
}

func (fsl *FileSizeLimiter) loadSettings() (err error) {
	var xmlSettings *settings.XmlSettings
	xmlSettings, err = settings.NewXmlSettings(fsl.settingsSourceFilePath)
	if err != nil {
		return err
	}

	fsl.fileSizeLimitPerMimeType = make(map[mimetype.Template]FileSize)

	for _, mimeType := range xmlSettings.FileSizeLimiter.MimeType {
		mimeTypeName := mimetype.Template(mimeType.Name)

		if !fsl.isMimeTypeAllowed(mimeTypeName) {
			continue
		}

		_, recordAlreadyExists := fsl.fileSizeLimitPerMimeType[mimeTypeName]
		if recordAlreadyExists {
			return errors.Errorf(ErrFDuplicateMimeType, mimeTypeName)
		}

		if mimeType.SizeLimit > InternalFileSizeLimit {
			return errors.Errorf(ErrFFileSizeDidNotPassInternalCheck, mimeTypeName)
		}

		fsl.fileSizeLimitPerMimeType[mimeTypeName] = FileSize(mimeType.SizeLimit)
	}

	return nil
}

func (fsl *FileSizeLimiter) isMimeTypeAllowed(mimeType mimetype.Template) (mimeTypeIsAllowed bool) {
	_, err := fileext.GetFileExtensions(mimeType)
	if err != nil {
		return false
	}

	return true
}

func (fsl *FileSizeLimiter) GetFileSizeLimit(mimeType mimetype.Template) (sizeLimit uint, err error) {
	var ok bool
	sizeLimit, ok = fsl.fileSizeLimitPerMimeType[mimeType]
	if !ok {
		return 0, errors.Errorf(ErrFMimeTypeNotAvailable, mimeType)
	}

	return sizeLimit, nil
}
