package message

// Форматы сообщений.
const (
	// MsgFSignalReceived -- сообщение "принят сигнал".
	MsgFSignalReceived = "signal %v is received"
)

// Сообщения.
const (
	// MsgKafkaNotReady -- сообщение "Kafka не готов(а)".
	MsgKafkaNotReady = "kafka is not ready"

	// MsgKafkaStart -- сообщение "потребитель Kafka стартовал".
	MsgKafkaStart = "kafka consumer has started"

	// MsgKafkaStop -- сообщение "потребитель Kafka остановился".
	MsgKafkaStop = "kafka consumer has stopped"

	// MsgSomethingNotReady -- сообщение "что-то не готово".
	MsgSomethingNotReady = "something is not ready"

	// MsgHttpServerStarting -- сообщение "запуск HTTP сервера".
	MsgHttpServerStarting = "http server is starting ..."

	// MsgHttpServerStopped -- сообщение "HTTP сервер остановлен".
	MsgHttpServerStopped = "http server is stopped"

	// MsgHttpServerError -- сообщение "ошибка HTTP сервера".
	MsgHttpServerError = "http server error"

	// MsgCriticalError -- сообщение "критическая ошибка".
	MsgCriticalError = "critical error"

	// MsgPrefixSeparator -- разделитель для сообщений, ставится после префикса.
	MsgPrefixSeparator = " "

	// MsgPrefixBusinessLogics -- префикс для сообщений -- "бизнес логика".
	MsgPrefixBusinessLogics = "business logics"

	// MsgPrefixMetrics -- префикс для сообщений -- "метрики".
	MsgPrefixMetrics = "metrics"
)

// ComposeMessageWithPrefix создаёт сообщение с префиксом.
func ComposeMessageWithPrefix(prefix string, msg string) string {
	return prefix + MsgPrefixSeparator + msg
}
