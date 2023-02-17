package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/kr/pretty"
	"github.com/rs/zerolog"
	"github.com/vault-thirteen/junk/gfe/internal/application/config"
	"github.com/vault-thirteen/junk/gfe/internal/message"
	"github.com/vault-thirteen/junk/gfe/internal/prometheus"
	storageInterface "github.com/vault-thirteen/junk/gfe/pkg/repository"
	"go.uber.org/atomic"
)

// ComponentsCount -- количество внутренних компонентов Kafka.
// В данный момент Kafka содержит следующие компоненты:
//  1. потребитель (читатель).
const ComponentsCount = 1

// MsgFDebugConfig -- формат сообщения для отладки настроек.
const MsgFDebugConfig = "event reader configuration: %+v"

// Kafka -- инфраструктура Kafka.
type Kafka struct {
	// Внешние объекты.
	// Обнуление и изменение этих объектов запрещено.
	logger            *zerolog.Logger
	storage           storageInterface.Storage
	prometheusMetrics *prometheus.Metrics

	eventReaderConfig    *config.MessageReader
	saramaConfig         *sarama.Config
	consumerGroupHandler *ConsumerGroupHandler
	consumerGroup        sarama.ConsumerGroup
	consumerWG           *sync.WaitGroup

	// Канал для оповещения внутренних компонентов о надобности завершения.
	// Компоненты читают из этого канала сигналы о завершении.
	// Следовательно, для завершения всех компонентов, нужно отправлять N
	// сигналов, где N равно количеству компонентов.
	close chan bool

	// Количество внутренних компонентов Kafka.
	// Этот параметр важен для правильной остановки внутренних компонентов.
	componentsCount byte
}

// NewKafka -- конструктор инфраструктуры Kafka.
func NewKafka(
	logger *zerolog.Logger,
	storage storageInterface.Storage,
	prometheusMetrics *prometheus.Metrics,
) (*Kafka, error) {
	k := new(Kafka)

	// Сохранение указателей на внешние объекты.
	k.logger = logger
	k.storage = storage
	k.prometheusMetrics = prometheusMetrics

	err := k.init()
	if err != nil {
		return nil, err
	}

	k.logger.Debug().Msg(pretty.Sprintf(MsgFDebugConfig, k.eventReaderConfig))

	return k, nil
}

// init производит первичную настройку инфраструктуры Kafka.
func (k *Kafka) init() (err error) {
	k.componentsCount = ComponentsCount
	k.close = make(chan bool, k.componentsCount)

	k.eventReaderConfig, err = config.GetEventReaderConfig()
	if err != nil {
		return err
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = config.SaramaClientID
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	k.saramaConfig = saramaConfig

	k.consumerGroupHandler = &ConsumerGroupHandler{
		isReady: new(atomic.Bool),
		parent:  k,
	}
	k.consumerGroupHandler.isReady.Store(false)

	k.consumerGroup, err = sarama.NewConsumerGroup(
		k.eventReaderConfig.KafkaBrokerAddressList,
		k.eventReaderConfig.KafkaConsumerGroupID,
		k.saramaConfig,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetConsumerGroupHandler возвращает обработчик группы потребителей.
func (k *Kafka) GetConsumerGroupHandler() *ConsumerGroupHandler {
	return k.consumerGroupHandler
}

// Wait ждёт готовности Kafka.
func (k *Kafka) Wait() {
	for {
		if k.consumerGroupHandler.isReady.Load() {
			break
		} else {
			time.Sleep(time.Millisecond * config.KafkaReadinessWaitIntervalMs)
		}
	}
}

// Run -- основной рабочий процесс Kafka.
func (k *Kafka) Run() {
	k.consumerWG = new(sync.WaitGroup)
	k.consumerWG.Add(1)

	go k.consumeMessages()
}

// consumeMessages читает сообщения из Kafka и обрабатывает их.
func (k *Kafka) consumeMessages() {
	defer k.consumerWG.Done()

	k.logger.Info().Msg(message.MsgKafkaStart)

	var err error

	for {
		if k.isShutdownRequired() {
			break
		}

		err = k.consumerGroup.Consume(
			context.Background(),
			k.eventReaderConfig.KafkaTopicList,
			k.consumerGroupHandler,
		)
		if err != nil {
			k.prometheusMetrics.KafkaConsumerGroupErrorsCount.Inc()
			k.logger.Err(err).Send()
			time.Sleep(time.Second * config.ConsumeDelaySecAfterError)
		}

		k.setReadinessState(false)
	}

	k.logger.Info().Msg(message.MsgKafkaStop)
}

// isShutdownRequired говорит, нужно ли останавливать Kafka.
// Если компонент, выполнивший этот метод, получил ответ 'true', то он теряет
// право опрашивать этот метод, поскольку, важно, чтобы все компоненты
// Kafka смогли получить сигнал завершения.
func (k *Kafka) isShutdownRequired() bool {
	select {
	case <-k.close:
		return true
	default:
		return false
	}
}

// Close останавливает Kafka.
func (k *Kafka) Close() (err error) {
	// Говорим всем компонентам, что нужно остановиться.
	var i byte
	for i = 1; i <= k.componentsCount; i++ {
		k.close <- true
	}

	err = k.consumerGroup.Close()
	if err != nil {
		return err
	}

	// Ждём завершения всех компонентов.
	k.consumerWG.Wait()

	return nil
}

// GetReadinessState возвращает состояние готовности.
func (k *Kafka) GetReadinessState() bool {
	return k.consumerGroupHandler.isReady.Load()
}

// setReadinessState устанавливает новое состояние готовности.
func (k *Kafka) setReadinessState(newState bool) {
	k.consumerGroupHandler.isReady.Store(newState)
}
