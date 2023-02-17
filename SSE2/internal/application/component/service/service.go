package service

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/vault-thirteen/junk/SSE2/internal/application/component/fsl"
	prometheusComponent "github.com/vault-thirteen/junk/SSE2/internal/application/component/metrics"
	workerComponent "github.com/vault-thirteen/junk/SSE2/internal/application/component/worker"
	"github.com/vault-thirteen/junk/SSE2/internal/application/config"
	"github.com/vault-thirteen/junk/SSE2/internal/messages"
	inputKafkaInterface "github.com/vault-thirteen/junk/SSE2/pkg/interfaces/kafka/input"
	outputKafkaInterface "github.com/vault-thirteen/junk/SSE2/pkg/interfaces/kafka/output"
	serviceInterface "github.com/vault-thirteen/junk/SSE2/pkg/interfaces/service"
	serviceRequestMessage "github.com/vault-thirteen/junk/SSE2/pkg/models/message/service/request"
	message "github.com/vault-thirteen/junk/SSE2/pkg/models/message/service/response"
	"go.uber.org/atomic"
)

const ErrorsChannelSize = 1024

type Service struct {
	logger            *zerolog.Logger
	fileSizeLimiter   *fsl.FileSizeLimiter
	prometheusMetrics *prometheusComponent.Metrics

	config            *config.Service
	inputKafka        inputKafkaInterface.Kafka
	outputKafka       outputKafkaInterface.Kafka
	storageSettings   *config.Storage
	starterLock       sync.Mutex
	stopperLock       sync.Mutex
	isStarted         atomic.Bool
	tasksChannel      chan *serviceRequestMessage.RequestMessage
	tasksReceiverWG   *sync.WaitGroup
	tasksWG           *sync.WaitGroup
	pendingTasksCount *atomic.Uint64
	errorsChannel     chan error
	idleWorkers       chan *workerComponent.Worker
	workersWG         *sync.WaitGroup
	converterVersion  string
}

const MsgFDebugConfig = "service configuration: %+v"

var (
	ErrInputKafkaIsNotReady    = errors.New("input kafka is not ready")
	ErrLoggerNull              = errors.New("logger is not set")
	ErrInputKafkaNull          = errors.New("input kafka is not set")
	ErrOutputKafkaIsNotReady   = errors.New("output kafka is not ready")
	ErrOutputKafkaNull         = errors.New("output kafka is not set")
	ErrStorageIsNotReady       = errors.New("storage is not ready")
	ErrStorageSettingsNull     = errors.New("storage settings are not set")
	ErrPrometheusMetricsNull   = errors.New("prometheus metrics are not set")
	ErrFileSizeLimiterNull     = errors.New("file size limiter is not set")
	ErrIsAlreadyStarted        = errors.New("is already started")
	ErrIsNotStarted            = errors.New("is not started")
	ErrTaskIsNull              = errors.New("task is not set")
	ErrTaskMessageIsNull       = errors.New("task message is not set")
	ErrTaskReturnAddressIsNull = errors.New("task return address is not set")
	ErrConverterVersion        = errors.New("converter version is unknown")
)

func NewService() (service serviceInterface.Service, err error) {
	s := new(Service)

	s.config, err = config.GetServiceConfig()
	if err != nil {
		return nil, err
	}

	s.tasksReceiverWG = new(sync.WaitGroup)
	s.tasksWG = new(sync.WaitGroup)
	s.pendingTasksCount = new(atomic.Uint64)
	s.errorsChannel = make(chan error, ErrorsChannelSize)
	s.workersWG = new(sync.WaitGroup)

	return s, nil
}

func (s *Service) Configure(
	logger *zerolog.Logger,
	fileSizeLimiter *fsl.FileSizeLimiter,
	inputKafka inputKafkaInterface.Kafka,
	outputKafka outputKafkaInterface.Kafka,
	storageSettings *config.Storage,
	prometheusMetrics *prometheusComponent.Metrics,
) (err error) {
	if logger == nil {
		return ErrLoggerNull
	}
	s.logger = logger

	if fileSizeLimiter == nil {
		return ErrFileSizeLimiterNull
	}
	s.fileSizeLimiter = fileSizeLimiter

	if inputKafka == nil {
		return ErrInputKafkaNull
	}
	s.inputKafka = inputKafka

	if outputKafka == nil {
		return ErrOutputKafkaNull
	}
	s.outputKafka = outputKafka

	if storageSettings == nil {
		return ErrStorageSettingsNull
	}
	s.storageSettings = storageSettings

	if prometheusMetrics == nil {
		return ErrPrometheusMetricsNull
	}
	s.prometheusMetrics = prometheusMetrics

	err = s.initLibreOffice()
	if err != nil {
		return err
	}

	err = s.initWorkersPool()
	if err != nil {
		return err
	}

	err = s.initConverter()
	if err != nil {
		return err
	}

	s.logger.Debug().Msg(pretty.Sprintf(MsgFDebugConfig, s.config))

	return nil
}

func (s *Service) initLibreOffice() (err error) {
	const MsgFolderIsCreated = "folder '%v' has been created"

	if !s.config.UseLibreOfficeMultipleUserInstallations {
		return nil
	}

	folderPath := filepath.Join(
		s.storageSettings.S3LocalFilesFolder,
		config.LibreOfficeTemporaryFolder,
		config.LibreOfficeUserInstallationFolder,
	)
	err = os.MkdirAll(folderPath, 0777)
	if err != nil {
		return err
	}

	s.logger.Info().Msgf(MsgFolderIsCreated, folderPath)

	for i := 1; i <= int(s.config.WorkersCount); i++ {
		folderPath = filepath.Join(
			s.storageSettings.S3LocalFilesFolder,
			config.LibreOfficeTemporaryFolder,
			config.LibreOfficeUserInstallationFolder,
			strconv.Itoa(i),
		)
		err = os.MkdirAll(folderPath, 0777)
		if err != nil {
			return err
		}

		s.logger.Debug().Msgf(MsgFolderIsCreated, folderPath)
	}

	return nil
}

func (s *Service) initWorkersPool() (err error) {
	s.idleWorkers = make(chan *workerComponent.Worker, s.config.WorkersCount)

	var worker *workerComponent.Worker
	for i := 1; i <= int(s.config.WorkersCount); i++ {
		worker, err = workerComponent.NewWorker(
			s.logger,
			s.fileSizeLimiter,
			s.prometheusMetrics,
			s.storageSettings,
			s.config,
			s.workersWG,
			uint(i),
		)
		if err != nil {
			return err
		}

		s.idleWorkers <- worker
	}

	return nil
}

func (s *Service) initConverter() (err error) {
	cmd := s.config.PathToConverterExecutable

	var converterHelpRawOutput []byte
	converterHelpRawOutput, err = exec.Command(cmd, "-h").CombinedOutput()
	if err != nil {
		return err
	}

	buffer := strings.Split(string(converterHelpRawOutput), "\n")
	if len(buffer) < 1 {
		return ErrConverterVersion
	}

	s.converterVersion = buffer[0]

	s.logger.Info().Msg("Converter: " + s.converterVersion)

	return nil
}

func (s *Service) Start() (err error) {
	s.starterLock.Lock()
	defer s.starterLock.Unlock()
	if s.isStarted.Load() {
		return ErrIsAlreadyStarted
	}

	s.tasksChannel, err = s.inputKafka.GetMessagesChannel()
	if err != nil {
		return err
	}

	s.tasksReceiverWG.Add(1)
	go s.receiveAndProcessTasks()

	s.isStarted.Store(true)

	return nil
}

func (s *Service) receiveAndProcessTasks() {
	defer s.tasksReceiverWG.Done()

	s.logger.Info().Msg(messages.MsgTasksReceiverStart)

	var task *serviceRequestMessage.RequestMessage
	for task = range s.tasksChannel {
		s.tasksWG.Add(1)
		s.pendingTasksCount.Inc()
		s.prometheusMetrics.PendingTasksCount.Inc()

		go s.receiveAndProcessTask(task)
	}

	s.tasksWG.Wait()

	s.logger.Info().Msg(messages.MsgTasksReceiverStop)
}

func (s *Service) receiveAndProcessTask(task *serviceRequestMessage.RequestMessage) {
	defer func() {
		s.prometheusMetrics.PendingTasksCount.Dec()
		s.pendingTasksCount.Dec()
		s.tasksWG.Done()
	}()

	s.logger.Debug().Msg(messages.MsgTaskReceiverStart)

	var response *message.ResponseMessage
	var err error
	response, err = s.processTask(task)
	if err != nil {
		s.errorsChannel <- err
	}

	if task.ReturnAddress != nil {
		task.ReturnAddress <- response
	}

	var taskMimeType string
	if (task != nil) && (task.KafkaMessage != nil) {
		taskMimeType = string(task.KafkaMessage.MimeType)
	}

	s.prometheusMetrics.ProcessedConversionRequestsCount.
		With(prometheus.Labels{
			config.MetricsLabelMimeType: taskMimeType,
		}).Inc()

	s.logger.Debug().Msg(messages.MsgTaskReceiverStop)
}

func (s *Service) processTask(
	task *serviceRequestMessage.RequestMessage,
) (result *message.ResponseMessage, err error) {
	if task == nil {
		return nil, ErrTaskIsNull
	}

	if task.KafkaMessage == nil {
		return nil, ErrTaskMessageIsNull
	}

	if task.ReturnAddress == nil {
		return nil, ErrTaskReturnAddressIsNull
	}

	timeWorkStart := time.Now()
	var workerNumber uint

	defer func() {
		result.ConversionResult.WorkerNumber = workerNumber

		timeWorkDuration := time.Since(timeWorkStart)

		result.ConversionResult.WorkTimeAsyncMs = uint(timeWorkDuration.Milliseconds())

		s.prometheusMetrics.ConversionDurationAsync.
			With(prometheus.Labels{
				config.MetricsLabelMimeType: string(task.KafkaMessage.MimeType),
			}).Observe(float64(timeWorkDuration.Milliseconds()))
	}()

	var worker = <-s.idleWorkers
	defer func() {
		s.idleWorkers <- worker
	}()

	workerNumber = worker.GetNumber()

	result, err = worker.ProcessTask(task)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *Service) Stop() (err error) {
	s.stopperLock.Lock()
	defer s.stopperLock.Unlock()
	if !s.isStarted.Load() {
		return ErrIsNotStarted
	}

	s.tasksReceiverWG.Wait()
	s.workersWG.Wait()

	close(s.errorsChannel)

	s.isStarted.Store(false)

	return nil
}

func (s *Service) GetErrorsChannel() chan error {
	return s.errorsChannel
}

func (s *Service) GetReadinessState() (isReady bool, err error) {
	if !s.inputKafka.GetReadinessState() {
		return false, ErrInputKafkaIsNotReady
	}

	if !s.outputKafka.GetReadinessState() {
		return false, ErrOutputKafkaIsNotReady
	}

	if !s.getStorageReadinessState() {
		return false, ErrStorageIsNotReady
	}

	return true, nil
}

func (s *Service) getStorageReadinessState() (isReady bool) {
	var worker = <-s.idleWorkers
	defer func() {
		s.idleWorkers <- worker
	}()

	if !worker.GetStorageReadinessState() {
		return false
	}

	return true
}
func (s *Service) SendConversionResults(message *message.ResponseMessage) (err error) {
	return s.outputKafka.SendConversionResults(message)
}

func (s *Service) GetFileSizeLimitSettingsFile() (filePath string) {
	return s.config.FileSizeLimitSettingsFile
}
