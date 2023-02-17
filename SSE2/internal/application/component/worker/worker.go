package worker

import (
	"context"
	"fmt"
	"image"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/vault-thirteen/junk/SSE2/internal/application/component/fsl"
	prometheusComponent "github.com/vault-thirteen/junk/SSE2/internal/application/component/metrics"
	"github.com/vault-thirteen/junk/SSE2/internal/application/component/storage"
	"github.com/vault-thirteen/junk/SSE2/internal/application/config"
	"github.com/vault-thirteen/junk/SSE2/internal/fileext"
	"github.com/vault-thirteen/junk/SSE2/internal/helper"
	storageInterface "github.com/vault-thirteen/junk/SSE2/pkg/interfaces/storage"
	"github.com/vault-thirteen/junk/SSE2/pkg/models/convertor/output"
	convertorResult "github.com/vault-thirteen/junk/SSE2/pkg/models/convertor/result"
	serviceRequestMessage "github.com/vault-thirteen/junk/SSE2/pkg/models/message/service/request"
	serviceResponseMessage "github.com/vault-thirteen/junk/SSE2/pkg/models/message/service/response"
	"github.com/vault-thirteen/junk/SSE2/pkg/models/mimetype"
	storageModels "github.com/vault-thirteen/junk/SSE2/pkg/models/storage"
	"go.uber.org/multierr"
)

const (
	FileBaseNameSuffixSmall           = "small"
	FileBaseNameSuffixLarge           = "large"
	MsgFLocalFileUploaded             = "local file '%v' has been uploaded to bucket '%v' at path '%v'"
	ExternalConverterExecutionTimeout = time.Minute * 5
)

type Worker struct {
	logger            *zerolog.Logger
	fileSizeLimiter   *fsl.FileSizeLimiter
	prometheusMetrics *prometheusComponent.Metrics
	serviceSettings   *config.Service
	storageSettings   *config.Storage

	storage                    storageInterface.Storage
	workersWG                  *sync.WaitGroup
	number                     uint
	pdfExtension               string
	pngExtension               string
	userInstallationFolderPath string
}

var (
	ErrLoggerNull             = errors.New("logger is not set")
	ErrFileSizeLimiterNull    = errors.New("file size limiter is not set")
	ErrWaitGroupNull          = errors.New("wait group is not set")
	ErrFFileSizeTooBig        = "file size (%v) is too big, limit is %v"
	ErrIntermediateResultNull = errors.New("intermediate result is not set")
	ErrFFileAlreadyExists     = "file '%v' already exists in bucket '%v'"
	ErrPrometheusMetricsNull  = errors.New("prometheus metrics are not set")
	ErrServiceSettingsNull    = errors.New("service settings are not set")
	ErrStorageSettingsNull    = errors.New("storage settings are not set")
)

func NewWorker(
	logger *zerolog.Logger,
	fileSizeLimiter *fsl.FileSizeLimiter,
	prometheusMetrics *prometheusComponent.Metrics,
	storageSettings *config.Storage,
	serviceSettings *config.Service,
	workersWG *sync.WaitGroup,
	number uint,
) (worker *Worker, err error) {
	w := new(Worker)

	if logger == nil {
		return nil, ErrLoggerNull
	}
	w.logger = logger

	if fileSizeLimiter == nil {
		return nil, ErrFileSizeLimiterNull
	}
	w.fileSizeLimiter = fileSizeLimiter

	if prometheusMetrics == nil {
		return nil, ErrPrometheusMetricsNull
	}
	w.prometheusMetrics = prometheusMetrics

	if serviceSettings == nil {
		return nil, ErrServiceSettingsNull
	}
	w.serviceSettings = serviceSettings

	if storageSettings == nil {
		return nil, ErrStorageSettingsNull
	}
	w.storageSettings = storageSettings

	w.storage, err = storage.NewStorage(w.logger, w.storageSettings)
	if err != nil {
		return nil, err
	}

	if workersWG == nil {
		return nil, ErrWaitGroupNull
	}
	w.workersWG = workersWG

	w.number = number

	w.pdfExtension, err = fileext.GetFileExtension(mimetype.ApplicationPdf)
	if err != nil {
		return nil, err
	}

	w.pngExtension, err = fileext.GetFileExtension(mimetype.ImagePng)
	if err != nil {
		return nil, err
	}

	if w.serviceSettings.UseLibreOfficeMultipleUserInstallations {
		w.userInstallationFolderPath = filepath.Join(
			w.storageSettings.S3LocalFilesFolder,
			config.LibreOfficeTemporaryFolder,
			config.LibreOfficeUserInstallationFolder,
			strconv.FormatUint(uint64(w.number), 10),
		)
	}

	return w, nil
}

func (w *Worker) ProcessTask(
	task *serviceRequestMessage.RequestMessage,
) (result *serviceResponseMessage.ResponseMessage, err error) {
	w.workersWG.Add(1)
	defer w.workersWG.Done()

	w.prometheusMetrics.IncomingConversionRequestsCount.
		With(prometheus.Labels{
			config.MetricsLabelMimeType: string(task.KafkaMessage.MimeType),
		}).Inc()

	timeWorkStart := time.Now()

	result = &serviceResponseMessage.ResponseMessage{
		KafkaMessage:     task.KafkaMessage,
		ConversionResult: &convertorResult.ConversionResult{},
	}

	defer func() {
		if err != nil {
			result.ConversionResult.Error = err
		}

		timeWorkDuration := time.Since(timeWorkStart)

		result.ConversionResult.WorkTimeByWorkerMs = uint(timeWorkDuration.Milliseconds())

		w.prometheusMetrics.ConversionDurationByWorker.
			With(prometheus.Labels{
				config.MetricsLabelMimeType: string(task.KafkaMessage.MimeType),
			}).Observe(float64(timeWorkDuration.Milliseconds()))
	}()

	err = w.checkTaskFileSize(task)
	if err != nil {
		return result, err
	}

	result.ConversionResult.LocalTemporaryFolderName = helper.MakeTemporaryFolderName()
	result.ConversionResult.LocalTemporaryFolderPath, err = helper.CreateSubFolder(
		w.storage.GetS3LocalFilesFolder(),
		result.ConversionResult.LocalTemporaryFolderName,
	)
	if err != nil {
		return result, err
	}

	defer func() {
		derr := w.deleteTemporaryData(result)
		if derr != nil {
			err = multierr.Combine(err, derr)
		}
	}()

	err = w.downloadAndSaveTaskFile(result, task)
	if err != nil {
		return result, err
	}

	err = w.convertTaskFileToPdf(result)
	if err != nil {
		return result, err
	}

	err = w.convertPdfFileFirstPageToPng(result)
	if err != nil {
		return result, err
	}

	err = w.makeTwoPngFiles(result)
	if err != nil {
		return result, err
	}

	err = w.uploadConvertedFiles(result, task)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (w *Worker) checkTaskFileSize(
	task *serviceRequestMessage.RequestMessage,
) (err error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), config.StorageFileSizeTimeout)
	defer cancelFn()

	var fileSize int
	fileSize, err = w.storage.GetFileSize(
		ctx,
		task.KafkaMessage.Bucket,
		task.KafkaMessage.FilePath,
	)
	if err != nil {
		return err
	}

	w.logger.Debug().Msgf("task=%+v, file size is %v", task, fileSize)

	var fileSizeLimit uint
	fileSizeLimit, err = w.fileSizeLimiter.GetFileSizeLimit(task.KafkaMessage.MimeType)
	if err != nil {
		return err
	}

	if uint(fileSize) > fileSizeLimit {
		return errors.Errorf(ErrFFileSizeTooBig, fileSize, fileSizeLimit)
	}

	return nil
}

func (w *Worker) downloadAndSaveTaskFile(
	intermediateResult *serviceResponseMessage.ResponseMessage,
	task *serviceRequestMessage.RequestMessage,
) (err error) {
	if intermediateResult == nil {
		return ErrIntermediateResultNull
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), config.StorageFileDownloadTimeout)
	defer cancelFn()

	var downloadResult *storageModels.DownloadResult
	downloadResult, err = w.storage.DownloadFile(
		ctx,
		task.KafkaMessage.Bucket,
		task.KafkaMessage.FilePath,
		intermediateResult.ConversionResult.LocalTemporaryFolderPath,
	)
	if err != nil {
		return err
	}

	intermediateResult.ConversionResult.LocalSourceFileName = downloadResult.LocalFileName
	intermediateResult.ConversionResult.LocalSourceFilePath = downloadResult.LocalFilePath

	return nil
}

func (w *Worker) convertTaskFileToPdf(
	intermediateResult *serviceResponseMessage.ResponseMessage,
) (err error) {
	if intermediateResult == nil {
		return ErrIntermediateResultNull
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), ExternalConverterExecutionTimeout)
	defer cancelFn()

	arguments := make([]string, 0, 7)

	if w.serviceSettings.UseLibreOfficeMultipleUserInstallations {
		userInstallationArgument := fmt.Sprintf(LibreOfficeArgumentFmtUserInstallation, w.userInstallationFolderPath)
		arguments = append(arguments, userInstallationArgument)
	}

	arguments = append(arguments,
		LibreOfficeArgumentConvertTo,
		LibreOfficeArgumentConversionFormatPdf,
		intermediateResult.ConversionResult.LocalSourceFilePath,
		LibreOfficeArgumentOutputFolder,
		intermediateResult.ConversionResult.LocalTemporaryFolderPath,
		LibreOfficeArgumentHeadless,
	)

	var processId *int
	var outputLines []string
	processId, outputLines, err = helper.ExecuteCommandAndGetOutput(
		ctx,
		w.logger,
		w.serviceSettings.PathToConverterExecutable,
		arguments,
	)
	if err != nil {
		if processId != nil {
			w.logger.Warn().Msgf("output of the process with PID=%v: %v", *processId, outputLines)
		}

		return err
	}

	w.logger.Debug().Msgf("output of the process with PID=%v: %v", *processId, outputLines)

	var converterOutput *output.Output
	converterOutput, err = output.ParseLibreOfficeConverterOutput(outputLines)
	if err != nil {
		return err
	}

	intermediateResult.ConversionResult.LocalPdfFilePath = converterOutput.DestinationFilePath

	return nil
}

func (w *Worker) convertPdfFileFirstPageToPng(
	intermediateResult *serviceResponseMessage.ResponseMessage,
) (err error) {
	if intermediateResult == nil {
		return ErrIntermediateResultNull
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), ExternalConverterExecutionTimeout)
	defer cancelFn()

	arguments := make([]string, 0, 7)

	if w.serviceSettings.UseLibreOfficeMultipleUserInstallations {
		userInstallationArgument := fmt.Sprintf(LibreOfficeArgumentFmtUserInstallation, w.userInstallationFolderPath)
		arguments = append(arguments, userInstallationArgument)
	}

	arguments = append(arguments,
		LibreOfficeArgumentConvertTo,
		LibreOfficeArgumentConversionFormatPng,
		intermediateResult.ConversionResult.LocalPdfFilePath,
		LibreOfficeArgumentOutputFolder,
		intermediateResult.ConversionResult.LocalTemporaryFolderPath,
		LibreOfficeArgumentHeadless,
	)

	var processId *int
	var outputLines []string
	processId, outputLines, err = helper.ExecuteCommandAndGetOutput(
		ctx,
		w.logger,
		w.serviceSettings.PathToConverterExecutable,
		arguments,
	)
	if err != nil {
		return err
	}

	w.logger.Debug().Msgf("output of the process with PID=%v: %v", *processId, outputLines)

	var converterOutput *output.Output
	converterOutput, err = output.ParseLibreOfficeConverterOutput(outputLines)
	if err != nil {
		return err
	}

	intermediateResult.ConversionResult.LocalFullSizeFirstPageFilePath = converterOutput.DestinationFilePath
	intermediateResult.ConversionResult.LocalFullSizeFirstPageFileName = filepath.Base(
		intermediateResult.ConversionResult.LocalFullSizeFirstPageFilePath)

	return nil
}

func (w *Worker) makeTwoPngFiles(
	intermediateResult *serviceResponseMessage.ResponseMessage,
) (err error) {
	if intermediateResult == nil {
		return ErrIntermediateResultNull
	}

	var fullSizePngImage image.Image
	fullSizePngImage, err = helper.GetImage(intermediateResult.ConversionResult.LocalFullSizeFirstPageFilePath)
	if err != nil {
		return err
	}

	smallImage := w.scaleImage(fullSizePngImage, int(w.serviceSettings.SmallPngImageMaximumSideDimension))

	intermediateResult.ConversionResult.LocalSmallFirstPageFilePath = filepath.Join(
		intermediateResult.ConversionResult.LocalTemporaryFolderPath,
		helper.AppendSuffixToFileBaseName(
			intermediateResult.ConversionResult.LocalFullSizeFirstPageFileName,
			FileBaseNameSuffixSmall,
		),
	)

	err = helper.SaveImageAsPngFile(
		smallImage,
		intermediateResult.ConversionResult.LocalSmallFirstPageFilePath,
	)
	if err != nil {
		return err
	}

	largeImage := w.scaleImage(fullSizePngImage, int(w.serviceSettings.LargePngImageMaximumSideDimension))

	intermediateResult.ConversionResult.LocalLargeFirstPageFilePath = filepath.Join(
		intermediateResult.ConversionResult.LocalTemporaryFolderPath,
		helper.AppendSuffixToFileBaseName(
			intermediateResult.ConversionResult.LocalFullSizeFirstPageFileName,
			FileBaseNameSuffixLarge,
		),
	)

	err = helper.SaveImageAsPngFile(
		largeImage,
		intermediateResult.ConversionResult.LocalLargeFirstPageFilePath,
	)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) scaleImage(
	fullSizeImage image.Image,
	scaledImageMaximumSideDimension int,
) (scaledImage image.Image) {
	scaleFactor := helper.GetScaleFactorForMaxSide(
		fullSizeImage.Bounds(),
		scaledImageMaximumSideDimension,
	)

	scaledImageRectangle := image.Rect(
		0,
		0,
		int(math.Round(float64(fullSizeImage.Bounds().Size().X)*scaleFactor)),
		int(math.Round(float64(fullSizeImage.Bounds().Size().Y)*scaleFactor)),
	)

	return helper.ScaleImage(fullSizeImage, scaledImageRectangle)
}

func (w *Worker) uploadConvertedFiles(
	intermediateResult *serviceResponseMessage.ResponseMessage,
	task *serviceRequestMessage.RequestMessage,
) (err error) {
	intermediateResult.ConversionResult.ConvertedPdfFileS3Path =
		task.KafkaMessage.FilePath + fileext.Separator + w.pdfExtension
	intermediateResult.ConversionResult.ConvertedSmallPngFileS3Path =
		task.KafkaMessage.FilePath + fileext.Separator + FileBaseNameSuffixSmall +
			fileext.Separator + w.pngExtension
	intermediateResult.ConversionResult.ConvertedLargePngFileS3Path =
		task.KafkaMessage.FilePath + fileext.Separator + FileBaseNameSuffixLarge +
			fileext.Separator + w.pngExtension

	err = w.uploadConvertedPdfFile(intermediateResult, task)
	if err != nil {
		return err
	}

	err = w.uploadConvertedSmallPngFile(intermediateResult, task)
	if err != nil {
		return err
	}

	err = w.uploadConvertedLargePngFile(intermediateResult, task)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) uploadConvertedPdfFile(
	intermediateResult *serviceResponseMessage.ResponseMessage,
	task *serviceRequestMessage.RequestMessage,
) (err error) {
	return w.checkAndUploadFile(
		intermediateResult.ConversionResult.LocalPdfFilePath,
		mimetype.ApplicationPdf,
		task.KafkaMessage.Bucket,
		intermediateResult.ConversionResult.ConvertedPdfFileS3Path,
	)
}

func (w *Worker) uploadConvertedSmallPngFile(
	intermediateResult *serviceResponseMessage.ResponseMessage,
	task *serviceRequestMessage.RequestMessage,
) (err error) {
	return w.checkAndUploadFile(
		intermediateResult.ConversionResult.LocalSmallFirstPageFilePath,
		mimetype.ImagePng,
		task.KafkaMessage.Bucket,
		intermediateResult.ConversionResult.ConvertedSmallPngFileS3Path,
	)
}

func (w *Worker) uploadConvertedLargePngFile(
	intermediateResult *serviceResponseMessage.ResponseMessage,
	task *serviceRequestMessage.RequestMessage,
) (err error) {
	return w.checkAndUploadFile(
		intermediateResult.ConversionResult.LocalLargeFirstPageFilePath,
		mimetype.ImagePng,
		task.KafkaMessage.Bucket,
		intermediateResult.ConversionResult.ConvertedLargePngFileS3Path,
	)
}

func (w *Worker) checkAndUploadFile(
	srcLocalFilePath string,
	contentType string,
	dstBucket string,
	dstFilePath string,
) (err error) {
	ctxA, cancelFnA := context.WithTimeout(context.Background(), config.StorageFileExistsTimeout)
	defer cancelFnA()

	var fileAlreadyExists *bool
	fileAlreadyExists, err = w.storage.DoesFileExist(ctxA, dstBucket, dstFilePath)
	if err != nil {
		return err
	}

	if *fileAlreadyExists {
		return errors.Errorf(ErrFFileAlreadyExists, dstFilePath, dstBucket)
	}

	ctxB, cancelFnB := context.WithTimeout(context.Background(), config.StorageFilePartUploadTimeout)
	defer cancelFnB()

	err = w.storage.UploadFile(ctxB, srcLocalFilePath, contentType, dstBucket, dstFilePath)
	if err != nil {
		return err
	}

	w.logger.Debug().Msgf(MsgFLocalFileUploaded, srcLocalFilePath, dstBucket, dstFilePath)

	return nil
}

func (w *Worker) deleteTemporaryData(
	intermediateResult *serviceResponseMessage.ResponseMessage,
) (err error) {
	if len(intermediateResult.ConversionResult.LocalSourceFilePath) > 0 {
		err = os.Remove(intermediateResult.ConversionResult.LocalSourceFilePath)
		if err != nil {
			return err
		}
	}

	if len(intermediateResult.ConversionResult.LocalPdfFilePath) > 0 {
		err = os.Remove(intermediateResult.ConversionResult.LocalPdfFilePath)
		if err != nil {
			return err
		}
	}

	if len(intermediateResult.ConversionResult.LocalFullSizeFirstPageFilePath) > 0 {
		err = os.Remove(intermediateResult.ConversionResult.LocalFullSizeFirstPageFilePath)
		if err != nil {
			return err
		}
	}

	if len(intermediateResult.ConversionResult.LocalSmallFirstPageFilePath) > 0 {
		err = os.Remove(intermediateResult.ConversionResult.LocalSmallFirstPageFilePath)
		if err != nil {
			return err
		}
	}

	if len(intermediateResult.ConversionResult.LocalLargeFirstPageFilePath) > 0 {
		err = os.Remove(intermediateResult.ConversionResult.LocalLargeFirstPageFilePath)
		if err != nil {
			return err
		}
	}

	if len(intermediateResult.ConversionResult.LocalTemporaryFolderPath) > 0 {
		err = os.Remove(intermediateResult.ConversionResult.LocalTemporaryFolderPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Worker) GetStorageReadinessState() (isReady bool) {
	return w.storage.IsReady()
}

func (w *Worker) GetNumber() uint {
	return w.number
}
