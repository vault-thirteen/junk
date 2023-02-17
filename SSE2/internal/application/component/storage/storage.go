package storage

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/vault-thirteen/junk/SSE2/internal/application/config"
	"github.com/vault-thirteen/junk/SSE2/internal/helper"
	storageInterface "github.com/vault-thirteen/junk/SSE2/pkg/interfaces/storage"
	"github.com/vault-thirteen/junk/SSE2/pkg/models/storage"
	"go.uber.org/multierr"
)

const FileSizeNone = -1

const AwsErrorCodeNotFound = "NotFound"

type Storage struct {
	logger   *zerolog.Logger
	settings *config.Storage

	sessionOptions session.Options
	s3Session      *session.Session
	s3Api          *s3.S3
}

var (
	ErrLoggerNull                        = errors.New("logger is not set")
	ErrStorageSettingsNull               = errors.New("storage settings are not set")
	ErrHeadObjectOutputContentLengthNull = errors.New("head object output content length is null")
	ErrHeadObjectOutputNull              = errors.New("head object output is null")
	ErrFUploadedDataSizeMismatch         = "uploaded data size mismatch: file size is %v, but uploaded bytes count is %v"
	ErrFStrangeError                     = "storage returned a strange error: %v"
	ErrApiNotInitialized                 = errors.New("api is not initialized")
)

func NewStorage(
	logger *zerolog.Logger,
	storageSettings *config.Storage,
) (storage storageInterface.Storage, err error) {
	s := new(Storage)

	if logger == nil {
		return nil, ErrLoggerNull
	}
	s.logger = logger

	if storageSettings == nil {
		return nil, ErrStorageSettingsNull
	}
	s.settings = storageSettings

	err = s.init()
	if err != nil {
		return nil, err
	}

	s.Wait()

	return s, nil
}

func (s *Storage) init() (err error) {
	s.sessionOptions = session.Options{
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentials(
				s.settings.S3AccessKey,
				s.settings.S3Secret,
				s.settings.S3Token,
			),
			Endpoint:         helper.NewStringPointer(s.settings.S3ServerAddress),
			Region:           aws.String(s.settings.S3Region),
			DisableSSL:       aws.Bool(s.settings.S3DisableSsl),
			S3ForcePathStyle: helper.NewBoolPointer(s.settings.S3ForcePathStyle),
		},
	}

	s.s3Session, err = session.NewSessionWithOptions(s.sessionOptions)
	if err != nil {
		return err
	}

	s.s3Api = s3.New(s.s3Session)

	return nil
}

func (s *Storage) Ping() (err error) {
	if s.s3Api == nil {
		return ErrApiNotInitialized
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), config.StorageListBucketsTimeout)
	defer cancelFn()

	var result *s3.ListBucketsOutput
	result, err = s.s3Api.ListBucketsWithContext(ctx, new(s3.ListBucketsInput))
	if err != nil {
		return err
	}

	s.logger.Debug().Msgf("buckets list: %v", result.String())

	return nil
}

func (s *Storage) IsReady() bool {
	err := s.Ping()
	if err != nil {
		return false
	}

	return true
}

func (s *Storage) Wait() {
	for {
		if s.IsReady() {
			break
		} else {
			time.Sleep(time.Second * config.StorageReadinessWaitIntervalSec)
		}
	}
}

func (s *Storage) GetS3LocalFilesFolder() (folderPath string) {
	return s.settings.S3LocalFilesFolder
}

func (s *Storage) GetFileSize(
	ctx context.Context,
	bucket string,
	filePath string,
) (fileSize int, err error) {
	var headObjectInput = &s3.HeadObjectInput{
		Bucket: helper.NewStringPointer(bucket),
		Key:    helper.NewStringPointer(filePath),
	}

	var headObjectOutput *s3.HeadObjectOutput
	headObjectOutput, err = s.s3Api.HeadObjectWithContext(ctx, headObjectInput)
	if err != nil {
		return FileSizeNone, err
	}

	if headObjectOutput == nil {
		return FileSizeNone, ErrHeadObjectOutputNull
	}

	if headObjectOutput.ContentLength == nil {
		return FileSizeNone, ErrHeadObjectOutputContentLengthNull
	}

	return int(*headObjectOutput.ContentLength), nil
}

func (s *Storage) DownloadFile(
	ctx context.Context,
	srcBucket string,
	srcFilePath string,
	dstLocalFolderPath string,
) (result *storage.DownloadResult, err error) {
	var getObjectInput = &s3.GetObjectInput{
		Bucket: helper.NewStringPointer(srcBucket),
		Key:    helper.NewStringPointer(srcFilePath),
	}

	var getObjectOutput *s3.GetObjectOutput
	getObjectOutput, err = s.s3Api.GetObjectWithContext(ctx, getObjectInput)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := getObjectOutput.Body.Close()
		if derr != nil {
			err = multierr.Combine(err, derr)
		}
	}()

	result = new(storage.DownloadResult)

	result.LocalFileName = helper.MakeTemporaryLocalFileName(srcFilePath)
	result.LocalFilePath = filepath.Join(
		dstLocalFolderPath,
		result.LocalFileName,
	)

	var localFile *os.File
	localFile, err = os.Create(result.LocalFilePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := localFile.Close()
		if derr != nil {
			err = multierr.Combine(err, derr)
		}
	}()

	_, err = io.Copy(localFile, getObjectOutput.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Storage) UploadFile(
	ctx context.Context,
	srcLocalFilePath string,
	contentType string,
	dstBucket string,
	dstFilePath string,
) (err error) {
	const MultipartUploadPartSizeBytes = 64 * 1024 * 1024

	var localFile *os.File
	localFile, err = os.Open(srcLocalFilePath)
	if err != nil {
		return err
	}

	defer func() {
		derr := localFile.Close()
		if derr != nil {
			err = multierr.Combine(err, derr)
		}
	}()

	fileInfo, _ := localFile.Stat()
	fileSize := int(fileInfo.Size())
	var uploadedDataSize int

	multipartUploadInput := &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(dstBucket),
		Key:         aws.String(dstFilePath),
		ContentType: aws.String(contentType),
	}

	var multipartUploadOutput *s3.CreateMultipartUploadOutput
	multipartUploadOutput, err = s.s3Api.CreateMultipartUploadWithContext(ctx, multipartUploadInput)
	if err != nil {
		return err
	}

	upload := new(s3.CompletedMultipartUpload)

	s.logger.Debug().Msgf("multipart upload initialization parameters: %v", multipartUploadOutput.String())

	buffer := make([]byte, MultipartUploadPartSizeBytes)

	defer func() {
		if err == nil {
			return
		}

		s.logger.Warn().Msgf("aborting multipart upload with id %v", multipartUploadOutput.UploadId)

		abortRequest := &s3.AbortMultipartUploadInput{
			Bucket:   multipartUploadOutput.Bucket,
			Key:      multipartUploadOutput.Key,
			UploadId: multipartUploadOutput.UploadId,
		}

		var abortResult *s3.AbortMultipartUploadOutput
		var derr error
		abortResult, derr = s.s3Api.AbortMultipartUploadWithContext(ctx, abortRequest)
		if derr != nil {
			err = multierr.Combine(err, derr)
			return
		}

		s.logger.Warn().Msgf("multipart upload abort result: %v", abortResult.String())
	}()

	var readBytesCount int
	var uploadPartNumber = 1
	for {
		readBytesCount, err = localFile.Read(buffer)
		if err != nil {
			if (err == io.EOF) && (readBytesCount == 0) {
				break
			}

			return err
		}

		partUploadRequest := &s3.UploadPartInput{
			Body:          bytes.NewReader(buffer[:readBytesCount]),
			Bucket:        multipartUploadOutput.Bucket,
			Key:           multipartUploadOutput.Key,
			PartNumber:    aws.Int64(int64(uploadPartNumber)),
			UploadId:      multipartUploadOutput.UploadId,
			ContentLength: aws.Int64(int64(readBytesCount)),
		}

		var partUploadResult *s3.UploadPartOutput
		partUploadResult, err = s.s3Api.UploadPartWithContext(ctx, partUploadRequest)
		if err != nil {
			return err
		}

		s.logger.Debug().Msgf(
			"multipart part #%v upload result: %v",
			uploadPartNumber,
			partUploadResult.String(),
		)

		uploadedDataSize += readBytesCount

		completedPart := &s3.CompletedPart{
			ETag:       partUploadResult.ETag,
			PartNumber: aws.Int64(int64(uploadPartNumber)),
		}

		upload.Parts = append(upload.Parts, completedPart)

		uploadPartNumber++

		for i := range buffer {
			buffer[i] = 0
		}
	}

	if uploadedDataSize != fileSize {
		err = errors.Errorf(ErrFUploadedDataSizeMismatch, fileSize, uploadedDataSize)
		return err
	}

	requestToCompleteUpload := &s3.CompleteMultipartUploadInput{
		Bucket:          multipartUploadOutput.Bucket,
		Key:             multipartUploadOutput.Key,
		UploadId:        multipartUploadOutput.UploadId,
		MultipartUpload: upload,
	}

	var result *s3.CompleteMultipartUploadOutput
	result, err = s.s3Api.CompleteMultipartUploadWithContext(ctx, requestToCompleteUpload)
	if err != nil {
		return err
	}

	s.logger.Debug().Msgf("multipart upload result: %v", result.String())

	return nil
}

func (s *Storage) DoesFileExist(
	ctx context.Context,
	bucket string,
	filePath string,
) (fileExists *bool, err error) {
	var headObjectInput = &s3.HeadObjectInput{
		Bucket: helper.NewStringPointer(bucket),
		Key:    helper.NewStringPointer(filePath),
	}

	var headObjectOutput *s3.HeadObjectOutput
	headObjectOutput, err = s.s3Api.HeadObjectWithContext(ctx, headObjectInput)
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if !ok {
			return nil, errors.Errorf(ErrFStrangeError, err)
		}

		if awsErr.Code() == AwsErrorCodeNotFound {
			return helper.NewBoolPointer(false), nil
		}

		return nil, err
	}

	if headObjectOutput == nil {
		return nil, ErrHeadObjectOutputNull
	}

	return helper.NewBoolPointer(true), nil
}
