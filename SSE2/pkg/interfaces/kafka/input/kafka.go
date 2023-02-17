package kafka

import (
	"io"

	serviceRequestMessage "github.com/vault-thirteen/junk/SSE2/pkg/models/message/service/request"
)

type Kafka interface {
	io.Closer
	Run()
	Wait()
	GetReadinessState() bool
	GetMessagesChannel() (chan *serviceRequestMessage.RequestMessage, error)
}
