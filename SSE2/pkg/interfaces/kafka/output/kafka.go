package kafka

import (
	"io"

	message "github.com/vault-thirteen/junk/SSE2/pkg/models/message/service/response"
)

type Kafka interface {
	io.Closer
	GetReadinessState() bool
	SendConversionResults(*message.ResponseMessage) error
}
