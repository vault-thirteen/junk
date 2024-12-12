package settings

const (
	ExeMediaInfo             = "MediaInfo"
	ServerShutdownTimeoutSec = 15
)

var allowedVideoFileExtensions = map[string]bool{
	"avi":  true, // AVI -> "video/x-msvideo".
	"f4v":  true, // ??? -> "video/x-f4v".
	"m2ts": true, // MPEG-TS -> "video/mp2t".
	"m4v":  true, // ??? -> "video/x-m4v".
	"mkv":  true, // Matroska -> "video/x-matroska".
	"mov":  true, // ??? -> "video/quicktime".
	"mp4":  true, // ??? -> "video/mp4".
	"mpeg": true, // ??? -> "video/mpeg".
	"mpg":  true, // ??? -> "video/mpeg".
	"ts":   true, // MPEG-TS -> "video/mp2t".
	"wmv":  true, // ??? -> "video/x-ms-wmv".
}

const (
	QueryParameterId      = "id"
	QueryParameterMovieId = "movie_id"
	QueryParameterTagId   = "tag_id"
)

const HeaderContentType = "Content-Type"

const HeaderContentTypeJson = "application/json"

const TableNameFormat_TagNMovies = "Tag_%d_Movies"

const GuiUrlFolder = "gui"

// IsVideoFileExtensionAllowed checks whether the specified video file
// extension is allowed to be used.
func IsVideoFileExtensionAllowed(ext string) bool {
	_, ok := allowedVideoFileExtensions[ext]

	return ok
}
