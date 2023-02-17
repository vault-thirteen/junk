package message

import (
	"github.com/vault-thirteen/junk/SSE2/pkg/models/mimetype"
)

type RequestMessage struct {
	MimeType mimetype.Template `json:"mimeType"`
	Bucket   string            `json:"bucket"`
	FilePath string            `json:"filePath"`
}
