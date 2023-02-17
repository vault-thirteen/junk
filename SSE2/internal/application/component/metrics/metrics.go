package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	IncomingKafkaMessagesCount       prometheus.Counter
	IncomingConversionRequestsCount  *prometheus.CounterVec
	PendingTasksCount                prometheus.Gauge
	ProcessedConversionRequestsCount *prometheus.CounterVec
	ConversionDurationByWorker       *prometheus.HistogramVec
	ConversionDurationAsync          *prometheus.HistogramVec
	KafkaConsumerErrorsCount         prometheus.Counter
	KafkaConsumerGroupErrorsCount    prometheus.Counter
	KafkaProducerErrorsCount         prometheus.Counter
}
