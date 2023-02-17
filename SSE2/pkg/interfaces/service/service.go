package service

import (
	"github.com/rs/zerolog"
	"github.com/vault-thirteen/junk/SSE2/internal/application/component/fsl"
	"github.com/vault-thirteen/junk/SSE2/internal/application/component/metrics"
	"github.com/vault-thirteen/junk/SSE2/internal/application/config"
	inputKafkaInterface "github.com/vault-thirteen/junk/SSE2/pkg/interfaces/kafka/input"
	outputKafkaInterface "github.com/vault-thirteen/junk/SSE2/pkg/interfaces/kafka/output"
	message "github.com/vault-thirteen/junk/SSE2/pkg/models/message/service/response"
)

type Service interface {
	Configure(
		logger *zerolog.Logger,
		fileSizeLimiter *fsl.FileSizeLimiter,
		inputKafka inputKafkaInterface.Kafka,
		outputKafka outputKafkaInterface.Kafka,
		storageSettings *config.Storage,
		prometheusMetrics *metrics.Metrics,
	) (err error)

	Start() (err error)
	Stop() (err error)
	GetErrorsChannel() chan error
	GetReadinessState() (bool, error)
	SendConversionResults(message *message.ResponseMessage) error
	GetFileSizeLimitSettingsFile() (filePath string)
}
