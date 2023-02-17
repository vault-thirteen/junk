package config

import "github.com/pkg/errors"

var (
	ErrHost                              = errors.New("host is empty")
	ErrPort                              = errors.New("port is empty")
	ErrAddress                           = errors.New("address is empty")
	ErrAccessKey                         = errors.New("access key is empty")
	ErrSecret                            = errors.New("secret is empty")
	ErrRegion                            = errors.New("region is empty")
	ErrLocalFilesFolder                  = errors.New("local files folder is not set")
	ErrConsumerGroupID                   = errors.New("consumer group id is empty")
	ErrBrokerAddressListEmpty            = errors.New("broker address list is empty")
	ErrTopicListEmpty                    = errors.New("topic list is empty")
	ErrWorkersCount                      = errors.New("workers count is wrong")
	ErrPathToConverterExecutable         = errors.New("path to converter executable is not set")
	ErrLargePngImageMaximumSideDimension = errors.New("large png image maximum side dimension is not valid")
	ErrSmallPngImageMaximumSideDimension = errors.New("small png image maximum side dimension is not valid")
	ErrFileSizeLimitSettingsFile         = errors.New("file size limit settings file is not set")
)
