package kafka

import (
	"encoding/json"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/vault-thirteen/junk/SSE2/internal/application/component/metrics"
	"github.com/vault-thirteen/junk/SSE2/internal/application/config"
	"github.com/vault-thirteen/junk/SSE2/internal/helper"
	outputKafkaInterface "github.com/vault-thirteen/junk/SSE2/pkg/interfaces/kafka/output"
	kafkaResponseMessage "github.com/vault-thirteen/junk/SSE2/pkg/models/message/kafka/response"
	serviceResponseMessage "github.com/vault-thirteen/junk/SSE2/pkg/models/message/service/response"
)

const MsgFDebugConfig = "message writer configuration: %+v"

type Kafka struct {
	logger            *zerolog.Logger
	prometheusMetrics *metrics.Metrics

	messageWriterConfiguration *config.MessageWriter
	saramaConfig               *sarama.Config

	client   sarama.Client
	producer sarama.SyncProducer
}

var (
	ErrClientIsNull         = errors.New("client is null")
	ErrNoTopics             = errors.New("no topics")
	ErrServiceMessageNull   = errors.New("service message is not set")
	ErrKafkaMessageNull     = errors.New("kafka message is not set")
	ErrConversionResultNull = errors.New("conversion result is not set")
)

func NewKafka(
	logger *zerolog.Logger,
	prometheusMetrics *metrics.Metrics,
) (outputKafkaInterface.Kafka, error) {
	k := new(Kafka)

	k.logger = logger
	k.prometheusMetrics = prometheusMetrics

	err := k.init()
	if err != nil {
		return nil, err
	}

	k.logger.Debug().Msg(pretty.Sprintf(MsgFDebugConfig, k.messageWriterConfiguration))

	return k, nil
}

func (k *Kafka) init() (err error) {
	k.messageWriterConfiguration, err = config.GetMessageWriterConfig()
	if err != nil {
		return err
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = config.SaramaClientID
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true
	k.saramaConfig = saramaConfig

	k.client, err = sarama.NewClient(
		k.messageWriterConfiguration.KafkaBrokerAddressList,
		k.saramaConfig,
	)
	if err != nil {
		return err
	}

	k.producer, err = sarama.NewSyncProducerFromClient(k.client)
	if err != nil {
		return err
	}

	return nil
}

func (k *Kafka) Close() (err error) {
	return k.producer.Close()
}

func (k *Kafka) GetReadinessState() (isReady bool) {
	const LogWriter = "output kafka"

	if k.client == nil {
		k.logger.Err(ErrClientIsNull).Msg(LogWriter)
		return false
	}

	topics, err := k.client.Topics()
	if err != nil {
		k.logger.Err(err).Msg(LogWriter)
		return false
	}

	if len(topics) < 1 {
		k.logger.Err(ErrNoTopics).Msg(LogWriter)
		return false
	}

	return true
}

func (k *Kafka) SendConversionResults(serviceMessage *serviceResponseMessage.ResponseMessage) (err error) {
	defer func() {
		if err != nil {
			k.prometheusMetrics.KafkaProducerErrorsCount.Inc()
		}
	}()

	var kafkaMessage *kafkaResponseMessage.ResponseMessage
	kafkaMessage, err = k.prepareConversionResultsMessage(serviceMessage)
	if err != nil {
		return err
	}

	var buffer []byte
	buffer, err = json.Marshal(kafkaMessage)
	if err != nil {
		return err
	}

	producerMessageHeaders := []sarama.RecordHeader{
		{
			Key:   []byte(kafkaResponseMessage.ResponseMessageHeaderNameWorkerNumber),
			Value: []byte(strconv.FormatUint(uint64(serviceMessage.ConversionResult.WorkerNumber), 10)),
		},
		{
			Key:   []byte(kafkaResponseMessage.ResponseMessageHeaderNameWorkTimeByWorkerMs),
			Value: []byte(strconv.FormatUint(uint64(serviceMessage.ConversionResult.WorkTimeByWorkerMs), 10)),
		},
		{
			Key:   []byte(kafkaResponseMessage.ResponseMessageHeaderNameWorkTimeAsyncMs),
			Value: []byte(strconv.FormatUint(uint64(serviceMessage.ConversionResult.WorkTimeAsyncMs), 10)),
		},
	}

	for _, topic := range k.messageWriterConfiguration.KafkaTopicList {
		producerMessage := &sarama.ProducerMessage{
			Topic:   topic,
			Value:   sarama.ByteEncoder(buffer),
			Headers: producerMessageHeaders,
		}

		_, _, err = k.producer.SendMessage(producerMessage)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k *Kafka) prepareConversionResultsMessage(
	serviceMessage *serviceResponseMessage.ResponseMessage,
) (kafkaMessage *kafkaResponseMessage.ResponseMessage, err error) {
	if serviceMessage == nil {
		return nil, ErrServiceMessageNull
	}

	if serviceMessage.KafkaMessage == nil {
		return nil, ErrKafkaMessageNull
	}

	if serviceMessage.ConversionResult == nil {
		return nil, ErrConversionResultNull
	}

	kafkaMessage = &kafkaResponseMessage.ResponseMessage{
		Task: serviceMessage.KafkaMessage,
		Result: &kafkaResponseMessage.ResponseMessageResult{
			Bucket:           serviceMessage.KafkaMessage.Bucket,
			PdfFilePath:      serviceMessage.ConversionResult.ConvertedPdfFileS3Path,
			SmallPngFilePath: serviceMessage.ConversionResult.ConvertedSmallPngFileS3Path,
			LargePngFilePath: serviceMessage.ConversionResult.ConvertedLargePngFileS3Path,
		},
	}

	if serviceMessage.IsResultSuccess() {
		kafkaMessage.Result.IsSuccess = true
		kafkaMessage.Result.Error = nil

		return kafkaMessage, nil
	}

	kafkaMessage.Result.IsSuccess = false
	kafkaMessage.Result.Error = helper.NewStringPointer(
		serviceMessage.ConversionResult.Error.Error(),
	)

	return kafkaMessage, nil
}
