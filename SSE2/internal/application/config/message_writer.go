package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/vault-thirteen/junk/SSE2/internal/helper"
)

type MessageWriter struct {
	KafkaBrokerAddressList []string `split_words:"true"`
	KafkaTopicList         []string `split_words:"true"`
}

func NewMessageWriter(envPrefix string) (cfg *MessageWriter, err error) {
	cfg = new(MessageWriter)
	err = envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (mwc *MessageWriter) IsValid() (bool, error) {
	if len(mwc.KafkaBrokerAddressList) < 1 {
		return false, ErrBrokerAddressListEmpty
	}

	if len(mwc.KafkaTopicList) < 1 {
		return false, ErrTopicListEmpty
	}

	return true, nil
}

func GetMessageWriterConfig() (writerConfig *MessageWriter, err error) {
	envPrefix := helper.ConcatenateEnvVarPrefixes(
		EnvironmentVariablePrefixApplication,
		EnvironmentVariablePrefixKafkaOutput,
	)

	writerConfig, err = NewMessageWriter(envPrefix)
	if err != nil {
		return nil, err
	}

	_, err = writerConfig.IsValid()
	if err != nil {
		return nil, err
	}

	return writerConfig, nil
}
