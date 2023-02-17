package integration

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/vault-thirteen/junk/gfe/pkg/models/event"
	"github.com/vault-thirteen/junk/gfe/pkg/models/message"
)

// Вспомогательный объект для проведения интеграционного тестирования Kafka.
type Kafka struct {
	parent *Test

	// Настройки Kafka.
	QueueSettings *QueueSettings

	// Конфигурация для 'sarama'.
	SaramaConfig *sarama.Config

	// Продюсер для 'sarama'.
	Producer sarama.SyncProducer
}

// NewKafka -- конструктор Kafka.
func NewKafka(
	parent *Test,
	testId string,
) (k *Kafka, err error) {
	k = new(Kafka)

	k.parent = parent

	// QueueSettings.
	k.QueueSettings = &QueueSettings{
		KafkaTopic: composeTopicName(testId),
		KafkaBrokerAddressList: []string{
			"127.0.0.1:9093",
		},
	}

	// SaramaConfig.
	k.SaramaConfig = sarama.NewConfig()
	k.SaramaConfig.ClientID = composeId(IntegrationTestIdPrefix, "client", testId)
	k.SaramaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	k.SaramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	k.SaramaConfig.Producer.Return.Errors = true
	k.SaramaConfig.Producer.Return.Successes = true

	// Producer.
	k.Producer, err = sarama.NewSyncProducer(
		k.QueueSettings.KafkaBrokerAddressList,
		k.SaramaConfig,
	)
	if err != nil {
		return nil, err
	}

	return k, nil
}

// Stop останавливает Kafka.
func (k *Kafka) Stop() (err error) {
	err = k.Producer.Close()
	if err != nil {
		return err
	}

	return nil
}

// ProduceMessages создаёт и отправляет несколько сообщений в Kafka.
func (k *Kafka) ProduceMessages() (err error) {
	sentMessages := make([]*message.Message, len(test.CreatedEventTypes))
	var eventMessage *message.Message

	for i, eventType := range test.CreatedEventTypes {
		eventMessage, err = k.makeAndSendMessage(eventType, i+1)
		if err != nil {
			return err
		}

		sentMessages[i] = eventMessage

		// Эта пауза нужна для того, чтобы Kafka и база дынных успели
		// обработать сообщение.
		time.Sleep(time.Second * 1)
	}

	// Сохраняем результат в случае успеха.
	test.SentEventMessages = sentMessages

	return nil
}

// makeAndSendMessage создаёт сообщение и отправляет его в Kafka.
// Параметр 'n' -- номер создаваемого сообщения, используется для создания
// времени сообщения по таблице.
func (k *Kafka) makeAndSendMessage(eventType event.TypeID, n int) (eventMessage *message.Message, err error) {
	var e = k.createEventMessage(eventType, n)

	var msg *sarama.ProducerMessage
	msg, err = k.prepareEventMessage(e)
	if err != nil {
		return nil, err
	}

	var partition int32
	var offset int64
	partition, offset, err = k.Producer.SendMessage(msg)
	if err != nil {
		return nil, err
	}

	fmt.Printf(
		"a message has been sent to kafka; partition=%v, offset=%v\r\n",
		partition,
		offset,
	)

	return e, nil
}

// createEventMessage создаёт сообщение для отправки.
// Параметр 'n' -- номер создаваемого сообщения, используется для создания
// времени сообщения по таблице.
// Таблица времени используется для значительного упрощения проверки
// результатов, в том числе для упрощения расчёта агрегированных событий.
func (k *Kafka) createEventMessage(eventType event.TypeID, n int) (msg *message.Message) {
	return &message.Message{
		UserID:      UserId,
		FileID:      FileId,
		EventTypeID: eventType,
		EventTime:   k.parent.CreatedEventTimes[n-1],
	}
}

// prepareEventMessage подготавливает сообщение для Kafka.
func (k *Kafka) prepareEventMessage(eventMessage *message.Message) (msg *sarama.ProducerMessage, err error) {
	var buffer []byte
	buffer, err = json.Marshal(eventMessage)
	if err != nil {
		return nil, err
	}

	msg = &sarama.ProducerMessage{
		Topic:     k.QueueSettings.KafkaTopic,
		Partition: -1,
		Value:     sarama.ByteEncoder(buffer),
	}

	return msg, nil
}
