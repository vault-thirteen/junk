package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics -- набор показателей.
type Metrics struct {
	// #1.
	// RequestsCount -- метрика количества запросов к HTTP серверу бизнес
	// логики.
	RequestsCount *prometheus.CounterVec

	// #2.
	// KafkaConsumerErrorsCount -- метрика количества ошибок потребителя Kafka.
	KafkaConsumerErrorsCount prometheus.Counter

	// #3.
	// KafkaConsumerGroupErrorsCount -- метрика количества ошибок группы
	// потребителя Kafka.
	KafkaConsumerGroupErrorsCount prometheus.Counter

	// #4.
	// DatabaseConnectionErrorsCount -- метрика количества ошибок соединения
	// базы данных.
	DatabaseConnectionErrorsCount prometheus.Counter
}
