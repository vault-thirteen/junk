package settings

import (
	"errors"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
)

// CaptchaSettings are parameters of the captcha.
type CaptchaSettings struct {
	// StoreImages flag configures the output of created images. When it is set
	// to True, Captcha manager stores an image on a storage device and
	// responds only with the task ID. When set to False, Captcha manager does
	// not store an image on a storage device and responds with both image
	// binary data and task ID.
	StoreImages bool `json:"storeImages"`

	// ImagesFolder sets the storage folder for created images. This setting is
	// used together with the 'StoreImages' flag.
	ImagesFolder string `json:"imagesFolder"`

	// Dimensions of created images.
	ImageWidth  uint `json:"imageWidth"`
	ImageHeight uint `json:"imageHeight"`

	// Image's time to live, in seconds. Each image is deleted when this time
	// passes after its creation.
	ImageTtlSec uint `json:"imageTtlSec"`

	ClearImagesFolderAtStart bool `json:"clearImagesFolderAtStart"`

	// This setting allows to start an HTTP server for serving captcha saved
	// image files. This setting is used together with the 'StoreImages' flag.
	UseHttpServerForImages bool `json:"useHttpServerForImages"`

	// Parameters of the HTTP server serving saved image files.
	HttpServerHost string `json:"httpServerHost"`
	HttpServerPort uint16 `json:"httpServerPort"`
	HttpServerName string `json:"httpServerName"`

	// Cache.
	IsCachingEnabled bool `json:"isCachingEnabled"`
	CacheSizeLimit   int  `json:"cacheSizeLimit"`
	CacheVolumeLimit int  `json:"cacheVolumeLimit"`
	CacheRecordTtl   uint `json:"cacheRecordTtl"`
}

func (s CaptchaSettings) Check() (err error) {
	if s.StoreImages == true {
		if len(s.ImagesFolder) == 0 {
			return errors.New(c.MsgCaptchaServiceSettingError)
		}
	}
	if (s.ImageWidth == 0) ||
		(s.ImageHeight == 0) ||
		(s.ImageTtlSec == 0) {
		return errors.New(c.MsgCaptchaServiceSettingError)
	}

	if s.UseHttpServerForImages == true {
		if (len(s.HttpServerHost) == 0) ||
			(s.HttpServerPort == 0) ||
			(len(s.HttpServerName) == 0) {
			return errors.New(c.MsgCaptchaServiceSettingError)
		}
	}

	if s.IsCachingEnabled {
		if (s.CacheSizeLimit <= 0) ||
			(s.CacheVolumeLimit <= 0) ||
			(s.CacheRecordTtl == 0) {
			return errors.New(c.MsgCaptchaServiceSettingError)
		}
	}

	return nil
}
