package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/vault-thirteen/junk/gfe/internal/application/config"
	iMessage "github.com/vault-thirteen/junk/gfe/internal/message"
	"github.com/vault-thirteen/junk/gfe/pkg/models/message"
	"go.uber.org/atomic"
)

// ConsumerGroupHandler -- тип, реализующий интерфейс "sarama.ConsumerGroupHandler".
// https://github.com/Shopify/sarama/blob/0676fc297fca0c5dfba927c0c4f6925ce2ae1e77/consumer_group.go#L795
type ConsumerGroupHandler struct {
	// Флаг готовности. Сделан на основе примера от разработчиков библиотеки
	// sarama, https://github.com/Shopify/sarama/blob/master/examples/consumergroup/main.go.
	isReady *atomic.Bool

	// Указатель на родительский объект.
	parent *Kafka
}

// Setup -- реализация метода Setup интерфейса "sarama.ConsumerGroupHandler".
// https://github.com/Shopify/sarama/blob/0676fc297fca0c5dfba927c0c4f6925ce2ae1e77/consumer_group.go#L795
func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	h.parent.setReadinessState(true)

	return nil
}

// Cleanup -- реализация метода Cleanup интерфейса "sarama.ConsumerGroupHandler".
// https://github.com/Shopify/sarama/blob/0676fc297fca0c5dfba927c0c4f6925ce2ae1e77/consumer_group.go#L795
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim -- реализация метода ConsumeClaim интерфейса "sarama.ConsumerGroupHandler".
// https://github.com/Shopify/sarama/blob/0676fc297fca0c5dfba927c0c4f6925ce2ae1e77/consumer_group.go#L795
func (h *ConsumerGroupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) (err error) {
	defer func() {
		if err != nil {
			h.parent.prometheusMetrics.KafkaConsumerErrorsCount.Inc()
			h.parent.logger.Err(err).Msg(iMessage.MsgCriticalError)
			time.Sleep(time.Second * config.ConsumeDelaySecAfterError)
		}

		var e = recover()
		if e != nil {
			err = errors.Errorf("%v", e)
			h.parent.logger.Err(err).Send()
		}
	}()

	var eventMessage *message.Message
	for claimMessage := range claim.Messages() {
		msg := claimMessage

		h.parent.logger.Debug().
			Str("message_topic", msg.Topic).
			Time("message_timestamp", msg.Timestamp).
			Str("message_value", string(msg.Value)).
			Msg("new claim message")

		// Декодируем сообщение.
		eventMessage = new(message.Message)
		err = json.Unmarshal(msg.Value, eventMessage)
		if err != nil {
			return err
		}

		eventMessage.RawMessage = msg

		// Производим нужные действия в базе данных.
		h.parent.storage.Wait()

		err = h.saveEvent(eventMessage)
		if err != nil {
			return err
		}

		// Помечаем сообщение как прочитанное.
		session.MarkMessage(eventMessage.RawMessage, "")
	}

	return nil
}

func (h *ConsumerGroupHandler) saveEvent(eventMessage *message.Message) (err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), config.StorageQueryTimeoutSec*time.Second)
	defer cancelFunc()

	err = h.parent.storage.SaveEvent(ctx, eventMessage)
	if err != nil {
		return err
	}

	return nil
}
