package message

import (
	"github.com/pkg/errors"
	"github.com/vault-thirteen/junk/SSE2/pkg/models/convertor/result"
	kafkaRequestMessage "github.com/vault-thirteen/junk/SSE2/pkg/models/message/kafka/request"
)

type ResponseMessage struct {
	KafkaMessage     *kafkaRequestMessage.RequestMessage
	ConversionResult *result.ConversionResult
}

var ErrResultNotSet = errors.New("result is not set")

func (rm *ResponseMessage) IsResultSuccess() (isResultSuccess bool) {
	if rm.ConversionResult == nil {
		return false
	}

	return rm.ConversionResult.IsSuccess()
}

func (rm *ResponseMessage) GetError() (err error) {
	if rm.ConversionResult == nil {
		return ErrResultNotSet
	}

	return rm.ConversionResult.Error
}
