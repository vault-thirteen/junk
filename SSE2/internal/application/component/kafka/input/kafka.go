package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/vault-thirteen/junk/SSE2/internal/application/component/metrics"
	"github.com/vault-thirteen/junk/SSE2/internal/application/config"
	"github.com/vault-thirteen/junk/SSE2/internal/messages"
	inputKafkaInterface "github.com/vault-thirteen/junk/SSE2/pkg/interfaces/kafka/input"
	serviceInterface "github.com/vault-thirteen/junk/SSE2/pkg/interfaces/service"
	serviceRequestMessage "github.com/vault-thirteen/junk/SSE2/pkg/models/message/service/request"
)

const (
	ComponentsCount     = 1
	MessagesChannelSize = 1024
	MsgFDebugConfig     = "message reader configuration: %+v"
)

type Kafka struct {
	logger            *zerolog.Logger
	prometheusMetrics *metrics.Metrics
	service           serviceInterface.Service

	messageReaderConfiguration *config.MessageReader
	saramaConfig               *sarama.Config
	consumerGroupHandler       *ConsumerGroupHandler
	consumerGroup              sarama.ConsumerGroup
	consumerWG                 *sync.WaitGroup

	messages chan *serviceRequestMessage.RequestMessage

	mustClose       chan bool
	componentsCount byte
}

var (
	ErrMessagesChannelIsNull = errors.New("messages channel is null")
)

func NewKafka(
	logger *zerolog.Logger,
	prometheusMetrics *metrics.Metrics,
	service serviceInterface.Service,
) (inputKafkaInterface.Kafka, error) {
	k := new(Kafka)

	k.logger = logger
	k.prometheusMetrics = prometheusMetrics
	k.service = service

	err := k.init()
	if err != nil {
		return nil, err
	}

	k.logger.Debug().Msg(pretty.Sprintf(MsgFDebugConfig, k.messageReaderConfiguration))

	return k, nil
}

func (k *Kafka) init() (err error) {
	k.componentsCount = ComponentsCount
	k.mustClose = make(chan bool, k.componentsCount)

	k.messageReaderConfiguration, err = config.GetMessageReaderConfig()
	if err != nil {
		return err
	}

	k.messages = make(chan *serviceRequestMessage.RequestMessage, MessagesChannelSize)

	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = config.SaramaClientID
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	k.saramaConfig = saramaConfig

	k.consumerGroupHandler, err = NewConsumerGroupHandler(k)
	if err != nil {
		return err
	}

	k.consumerGroup, err = sarama.NewConsumerGroup(
		k.messageReaderConfiguration.KafkaBrokerAddressList,
		k.messageReaderConfiguration.KafkaConsumerGroupID,
		k.saramaConfig,
	)
	if err != nil {
		return err
	}

	return nil
}

func (k *Kafka) GetConsumerGroupHandler() *ConsumerGroupHandler {
	return k.consumerGroupHandler
}

func (k *Kafka) Wait() {
	for {
		if k.consumerGroupHandler.isReady.Load() {
			break
		} else {
			time.Sleep(time.Millisecond * config.KafkaReadinessWaitIntervalMs)
		}
	}
}

func (k *Kafka) Run() {
	k.consumerWG = new(sync.WaitGroup)
	k.consumerWG.Add(1)

	go k.consumeMessages()
}

func (k *Kafka) consumeMessages() {
	defer k.consumerWG.Done()

	k.logger.Info().Msg(messages.MsgKafkaConsumerStart)

	var err error

	for {
		if k.isShutdownRequired() {
			break
		}

		err = k.consumerGroup.Consume(
			context.Background(),
			k.messageReaderConfiguration.KafkaTopicList,
			k.consumerGroupHandler,
		)
		if err != nil {
			k.prometheusMetrics.KafkaConsumerGroupErrorsCount.Inc()
			k.logger.Err(err).Send()
			time.Sleep(time.Second * config.ConsumeDelaySecAfterError)
		}

		k.setReadinessState(false)
	}

	k.logger.Info().Msg(messages.MsgKafkaConsumerStop)
}

func (k *Kafka) isShutdownRequired() bool {
	select {
	case <-k.mustClose:
		return true
	default:
		return false
	}
}

func (k *Kafka) Close() (err error) {
	var i byte
	for i = 1; i <= k.componentsCount; i++ {
		k.mustClose <- true
	}

	err = k.consumerGroup.Close()
	if err != nil {
		return err
	}

	k.consumerWG.Wait()

	close(k.messages)

	return nil
}

func (k *Kafka) GetReadinessState() (isReady bool) {
	return k.consumerGroupHandler.isReady.Load()
}

func (k *Kafka) setReadinessState(newState bool) {
	k.consumerGroupHandler.isReady.Store(newState)
}

func (k *Kafka) GetMessagesChannel() (chan *serviceRequestMessage.RequestMessage, error) {
	if k.messages == nil {
		return nil, ErrMessagesChannelIsNull
	}

	return k.messages, nil
}
