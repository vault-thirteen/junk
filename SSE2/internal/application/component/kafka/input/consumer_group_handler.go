package kafka

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/vault-thirteen/junk/SSE2/internal/application/config"
	"github.com/vault-thirteen/junk/SSE2/internal/messages"
	message "github.com/vault-thirteen/junk/SSE2/pkg/models/message/kafka/request"
	serviceRequestMessage "github.com/vault-thirteen/junk/SSE2/pkg/models/message/service/request"
	serviceResponseMessage "github.com/vault-thirteen/junk/SSE2/pkg/models/message/service/response"
	"go.uber.org/atomic"
)

const ConsumeErrorsChannelSize = 16

// https://github.com/Shopify/sarama/blob/0676fc297fca0c5dfba927c0c4f6925ce2ae1e77/consumer_group.go#L795
type ConsumerGroupHandler struct {
	// https://github.com/Shopify/sarama/blob/master/examples/consumergroup/main.go.
	isReady *atomic.Bool

	parent *Kafka

	consumeErrors     chan error
	firstConsumeError chan error

	processorsWG       *sync.WaitGroup
	asyncErrorReaderWG *sync.WaitGroup
}

func NewConsumerGroupHandler(
	parent *Kafka,
) (handler *ConsumerGroupHandler, err error) {
	handler = &ConsumerGroupHandler{
		isReady:            new(atomic.Bool),
		parent:             parent,
		processorsWG:       new(sync.WaitGroup),
		asyncErrorReaderWG: new(sync.WaitGroup),
	}

	handler.isReady.Store(false)

	return handler, nil
}

// https://github.com/Shopify/sarama/blob/0676fc297fca0c5dfba927c0c4f6925ce2ae1e77/consumer_group.go#L795
func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	h.parent.setReadinessState(true)

	return nil
}

// https://github.com/Shopify/sarama/blob/0676fc297fca0c5dfba927c0c4f6925ce2ae1e77/consumer_group.go#L795
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// https://github.com/Shopify/sarama/blob/0676fc297fca0c5dfba927c0c4f6925ce2ae1e77/consumer_group.go#L795
func (h *ConsumerGroupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) (err error) {
	defer func() {
		if err != nil {
			h.parent.prometheusMetrics.KafkaConsumerErrorsCount.Inc()
			h.parent.logger.Err(err).Msg(messages.MsgCriticalError)
			time.Sleep(time.Second * config.ConsumeDelaySecAfterError)
		}

		var e = recover()
		if e != nil {
			err = errors.Errorf("%v", e)
			h.parent.logger.Err(err).Send()
		}
	}()

	h.consumeErrors = make(chan error, ConsumeErrorsChannelSize)
	h.firstConsumeError = make(chan error, 1)

	for claimMessage := range claim.Messages() {
		msg := claimMessage

		h.processorsWG.Add(1)
		go h.processMessage(session, msg)
	}

	h.asyncErrorReaderWG.Add(1)

	go func() {
		defer h.asyncErrorReaderWG.Done()

		var atLeastOneErrorHasOccurred bool

		for asyncErr := range h.consumeErrors {
			if asyncErr == nil {
				continue
			}

			h.parent.logger.Err(asyncErr).Msg("ConsumerGroupHandler.processMessage")

			if !atLeastOneErrorHasOccurred {
				atLeastOneErrorHasOccurred = true
				h.firstConsumeError <- asyncErr
			}
		}
	}()

	h.processorsWG.Wait()
	close(h.consumeErrors)
	h.asyncErrorReaderWG.Wait()
	close(h.firstConsumeError)

	select {
	case err = <-h.firstConsumeError:
		return err

	default:
		return nil
	}
}

func (h *ConsumerGroupHandler) processMessage(
	session sarama.ConsumerGroupSession,
	msg *sarama.ConsumerMessage,
) {
	var err error

	defer func() {
		if err != nil {
			h.consumeErrors <- err
		}

		h.processorsWG.Done()
	}()

	h.parent.prometheusMetrics.IncomingKafkaMessagesCount.Inc()

	h.parent.logger.Debug().
		Str("message_topic", msg.Topic).
		Time("message_timestamp", msg.Timestamp).
		Str("message_value", string(msg.Value)).
		Msg("new claim message")

	var taskMessage = new(message.RequestMessage)
	err = json.Unmarshal(msg.Value, taskMessage)
	if err != nil {
		return
	}

	returnAddress := make(chan *serviceResponseMessage.ResponseMessage, 1)

	h.parent.messages <- &serviceRequestMessage.RequestMessage{
		KafkaMessage:  taskMessage,
		ReturnAddress: returnAddress,
	}

	taskResult := <-returnAddress

	err = h.parent.service.SendConversionResults(taskResult)
	if err != nil {
		return
	}

	if !taskResult.IsResultSuccess() {
		err = taskResult.GetError()

		return
	}

	session.MarkMessage(msg, "")
}
