package config

// Различные настройки.
const (
	// SaramaClientID -- идентификатор клиента sarama.
	SaramaClientID = "gfe"

	// ConsumeDelaySecAfterError -- период времени (в секундах), в течение
	// которого группа потребителя Kafka делает паузу после неудачной попытки
	// соединения с Kafka.
	ConsumeDelaySecAfterError = 5

	// KafkaReadinessWaitIntervalMs -- период времени (в миллисекундах), в
	// течение которого происходит задержка перед следующей попыткой проверки
	// готовности читателя Kafka.
	KafkaReadinessWaitIntervalMs = 100
)
