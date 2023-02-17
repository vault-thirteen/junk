package mimetype

type Template string

// https://www.iana.org/assignments/media-types/media-types.xhtml.
const (
	ApplicationMicrosoftExcel               = "application/vnd.ms-excel"
	ApplicationMicrosoftWord                = "application/msword"
	ApplicationMicrosoftWordMacro           = "application/vnd.ms-word.document.macroEnabled.12"
	ApplicationOasisOpenDocumentText        = "application/vnd.oasis.opendocument.text"
	ApplicationOasisOpenDocumentSpreadsheet = "application/vnd.oasis.opendocument.spreadsheet"
	ApplicationOfficeOpenXmlDocument        = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	ApplicationOfficeOpenXmlWorkbook        = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	ApplicationPdf                          = "application/pdf"
	ApplicationPostscript                   = "application/postscript"
	ApplicationRtf                          = "application/rtf"
	ApplicationVndAppleNumbers              = "application/vnd.apple.numbers"
	ApplicationVndApplePages                = "application/vnd.apple.pages"
	ApplicationVndWordPerfect               = "application/vnd.wordperfect"
	ApplicationWordPerfect51                = "application/wordperfect5.1"
)

const (
	ImagePng = "image/png"
)

const (
	FontOtf = "font/otf"
	FontTtf = "font/ttf"
)

const (
	TextCsv = "text/csv"
	TextRtf = "text/rtf"
)
