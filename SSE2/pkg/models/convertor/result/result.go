package result

type ConversionResult struct {
	Error                          error
	LocalTemporaryFolderName       string
	LocalTemporaryFolderPath       string
	LocalSourceFileName            string
	LocalSourceFilePath            string
	LocalPdfFilePath               string
	LocalFullSizeFirstPageFileName string
	LocalFullSizeFirstPageFilePath string
	LocalLargeFirstPageFilePath    string
	LocalSmallFirstPageFilePath    string
	ConvertedPdfFileS3Path         string
	ConvertedSmallPngFileS3Path    string
	ConvertedLargePngFileS3Path    string
	WorkerNumber                   uint
	WorkTimeByWorkerMs             uint
	WorkTimeAsyncMs                uint
}

func (cr *ConversionResult) IsSuccess() (isSuccess bool) {
	if cr.Error != nil {
		return false
	}

	return true
}
