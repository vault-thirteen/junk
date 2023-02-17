package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/vault-thirteen/junk/SSE2/internal/helper"
)

type MessageReader struct {
	KafkaConsumerGroupID   string   `split_words:"true"`
	KafkaBrokerAddressList []string `split_words:"true"`
	KafkaTopicList         []string `split_words:"true"`
}

func NewMessageReader(envPrefix string) (cfg *MessageReader, err error) {
	cfg = new(MessageReader)
	err = envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

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

func GetMessageReaderConfig() (readerConfig *MessageReader, err error) {
	envPrefix := helper.ConcatenateEnvVarPrefixes(
		EnvironmentVariablePrefixApplication,
		EnvironmentVariablePrefixKafkaInput,
	)

	readerConfig, err = NewMessageReader(envPrefix)
	if err != nil {
		return nil, err
	}

	_, err = readerConfig.IsValid()
	if err != nil {
		return nil, err
	}

	return readerConfig, nil
}
