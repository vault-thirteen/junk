package message

import (
	kafkaRequestMessage "github.com/vault-thirteen/junk/SSE2/pkg/models/message/kafka/request"
	serviceResponseMessage "github.com/vault-thirteen/junk/SSE2/pkg/models/message/service/response"
)

type RequestMessage struct {
	KafkaMessage  *kafkaRequestMessage.RequestMessage
	ReturnAddress chan *serviceResponseMessage.ResponseMessage
}
