package config

import (
	"github.com/kelseyhightower/envconfig"
)

// MessageReader -- настройка читателя сообщений.
type MessageReader struct {
	// Идентификатор группы потребителя сообщений Kafka.
	KafkaConsumerGroupID string `split_words:"true"`

	// Список адресов посредников для Kafka.
	KafkaBrokerAddressList []string `split_words:"true"`

	// Список тем для Kafka.
	KafkaTopicList []string `split_words:"true"`
}

// NewMessageReader -- создатель настроек читателя сообщений.
func NewMessageReader(envPrefix string) (cfg *MessageReader, err error) {
	cfg = new(MessageReader)
	err = envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// IsValid проверяет годность настройки читателя сообщений.
func (mrc *MessageReader) IsValid() (bool, error) {
	if len(mrc.KafkaConsumerGroupID) < 1 {
		return false, ErrConsumerGroupID
	}

	if len(mrc.KafkaBrokerAddressList) < 1 {
		return false, ErrBrokerAddressListEmpty
	}

	if len(mrc.KafkaTopicList) < 1 {
		return false, ErrTopicListEmpty
	}

	return true, nil
}

// GetEventReaderConfig получает настройки читателя событий.
func GetEventReaderConfig() (eventReaderConfig *MessageReader, err error) {
	eventReaderConfig, err = NewMessageReader(environmentVariablePrefixApplication)
	if err != nil {
		return nil, err
	}

	_, err = eventReaderConfig.IsValid()
	if err != nil {
		return nil, err
	}

	return eventReaderConfig, nil
}
