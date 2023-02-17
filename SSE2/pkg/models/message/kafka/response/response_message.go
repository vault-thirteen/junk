package message

import (
	message "github.com/vault-thirteen/junk/SSE2/pkg/models/message/kafka/request"
)

type ResponseMessage struct {
	Task   *message.RequestMessage `json:"task"`
	Result *ResponseMessageResult
}
