package message

type ResponseMessageResult struct {
	IsSuccess        bool    `json:"isSuccess"`
	Error            *string `json:"error"`
	Bucket           string  `json:"bucket"`
	PdfFilePath      string  `json:"pdfFilePath"`
	SmallPngFilePath string  `json:"smallPngFilePath"`
	LargePngFilePath string  `json:"largePngFilePath"`
}
