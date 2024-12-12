package models

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"path/filepath"

	af "github.com/vault-thirteen/auxie/file"
)

// FrontEndFileData is auxiliary data for a front end static file.
type FrontEndFileData struct {
	UrlPath     cm.Path
	FilePath    cm.Path
	ContentType string
	CachedFile  []byte
}

func NewFrontEndFileData(frontEndPath cm.Path, fileName cm.Path, contentType string, frontendAssetsFolder cm.Path) (fefd FrontEndFileData, err error) {
	fefd = FrontEndFileData{
		UrlPath:     frontEndPath + fileName,
		FilePath:    cm.Path(filepath.Join(frontendAssetsFolder.ToString(), fileName.ToString())),
		ContentType: contentType,
	}

	fefd.CachedFile, err = af.GetFileContents(fefd.FilePath.ToString())
	if err != nil {
		return fefd, err
	}

	return fefd, nil
}
