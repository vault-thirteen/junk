package integration

type QueueSettings struct {
	// Список адресов посредников для Kafka.
	KafkaBrokerAddressList []string

	// Тема (топик) для отправляемых сообщений.
	KafkaTopic string
}

type EnvironmentVariable struct {
	Name  string
	Value string
}
