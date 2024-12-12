package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"path"
	"spdb_manager/cli"
	"spdb_manager/cli/printer"
	"spdb_manager/database"
	"spdb_manager/helper"
	mi "spdb_manager/media_info"
	"spdb_manager/models"
	"spdb_manager/settings"
	"strings"
	"time"

	"go.uber.org/multierr"
)

// importNewMoviesList imports data about new movies into the database.
func importNewMoviesList(
	cliParams *cli.Parameters,
) (err error) {
	if len(cliParams.VideosListUrl) == 0 {
		return errors.New("videos list url is not set")
	}

	var mediaInfoVersionText string
	mediaInfoVersionText, err = mi.GetExeVersion()
	if err != nil {
		return err
	}
	printer.PrintNotice(mediaInfoVersionText)

	var files []string
	files, err = helper.GetFilesList(cliParams.VideosListUrl)
	if err != nil {
		return err
	}
	printer.PrintNormalText(fmt.Sprintf("%d files to process.", len(files)))

	err = database.CheckConnection(cliParams.DbFile)
	if err != nil {
		return err
	}

	var dbo *database.DbForFilesImport
	dbo, err = database.OpenForFilesImport(cliParams.DbFile)
	if err != nil {
		return err
	}

	defer func() {
		var derr error
		if err != nil {
			derr = dbo.Abort()
			if derr != nil {
				err = multierr.Combine(err, derr)
			}
		} else {
			derr = dbo.Close()
			if derr != nil {
				err = multierr.Combine(err, derr)
			}
		}
	}()

	err = processNewMovieFiles(dbo, files, cliParams.FilesToSkip)
	if err != nil {
		return err
	}

	return nil
}

// processNewMovieFiles tries to get data about new movie files in the list and
// insert it into the database. Files having inappropriate data are ignored.
func processNewMovieFiles(
	dbo *database.DbForFilesImport,
	files []string,
	filesToSkip uint,
) (err error) {
	var n uint = 1
	var movie *models.MovieRaw
	for _, filePath := range files {
		if n <= filesToSkip {
			n++
			continue
		}

		printer.PrintNotice(fmt.Sprintf("%d. %s", n, filePath))

		var isWarning bool
		movie, err, isWarning = getNewMovieData(dbo, filePath)
		if err != nil {
			if isWarning {
				printer.PrintInfo(err.Error())
			} else {
				printer.PrintError(err.Error())
			}

			n++
			continue
		}

		err = dbo.InsertNewMovie(movie)
		if err != nil {
			return err
		}

		n++
		continue
	}

	return nil
}

// getNewMovieData gets various data about a new movie.
func getNewMovieData(
	dbo *database.DbForFilesImport,
	filePath string,
) (movie *models.MovieRaw, err error, isWarning bool) {
	var mediaInfoResponse *mi.Response
	mediaInfoResponse, err = mi.GetData(filePath)
	if err != nil {
		return nil, err, false
	}

	fileExtNominal := strings.TrimPrefix(path.Ext(filePath), ".")

	err, isWarning = preCheckNewMovieMediaInfo(mediaInfoResponse, fileExtNominal)
	if err != nil {
		return nil, err, isWarning
	}

	movie = &models.MovieRaw{
		FileExtension:  mediaInfoResponse.Media.Tracks[0].FileExtension,
		OverallBitRate: mediaInfoResponse.Media.Tracks[0].OverallBitRate,
		FrameWidth:     mediaInfoResponse.Media.Tracks[1].Width,
		FrameHeight:    mediaInfoResponse.Media.Tracks[1].Height,
		Duration:       mediaInfoResponse.Media.Tracks[0].Duration,
	}

	err = fetchNewMovieAdditionalData(movie, dbo, filePath, mediaInfoResponse)
	if err != nil {
		return nil, err, false
	}

	movie.TimeOfRegistration = time.Now().Unix()

	return movie, nil, false
}

// preCheckNewMovieMediaInfo makes a preliminary check of the media information
// of a new movie.
func preCheckNewMovieMediaInfo(
	mediaInfoResponse *mi.Response,
	fileExtNominal string,
) (err error, isWarning bool) {
	if !settings.IsVideoFileExtensionAllowed(fileExtNominal) {
		return fmt.Errorf("skipping an unsupported file extension: %s",
			fileExtNominal), true
	}

	if !settings.IsVideoFileExtensionAllowed(mediaInfoResponse.Media.Tracks[0].FileExtension) {
		return fmt.Errorf("unexpected file extension: %s",
			mediaInfoResponse.Media.Tracks[0].FileExtension), false
	}

	if mediaInfoResponse == nil {
		return errors.New("response is not set"), false
	}

	if len(mediaInfoResponse.Media.Tracks) < 3 {
		return fmt.Errorf("media tracks error: %d", len(mediaInfoResponse.Media.Tracks)), true
	}

	if mediaInfoResponse.Media.Tracks[0].VideoCount != 1 {
		return fmt.Errorf("video tracks count error, expected 1, got %d",
			mediaInfoResponse.Media.Tracks[0].VideoCount), true
	}

	if mediaInfoResponse.Media.Tracks[0].AudioCount != 1 {
		return fmt.Errorf("audio tracks count error, expected 1, got %d",
			mediaInfoResponse.Media.Tracks[0].AudioCount), true
	}

	return nil, false
}

// fetchNewMovieAdditionalData fetches additional data for a new movie from the
// database and saves it into the movie object.
func fetchNewMovieAdditionalData(
	movie *models.MovieRaw, // <- Data Receiver.
	dbo *database.DbForFilesImport,
	filePath string,
	mediaInfoResponse *mi.Response,
) (err error) {
	var fileUrl *url.URL
	fileUrl, err = url.Parse(filePath)
	if err != nil {
		return err
	}

	movie.FilePath = path.Dir(fileUrl.Path)
	movie.FileName = path.Base(fileUrl.Path)
	movie.Title = helper.GetFileTitle(movie.FileName)
	serverAddress := fmt.Sprintf("%s://%s:%s", fileUrl.Scheme, fileUrl.Hostname(), fileUrl.Port())

	err = fetchNewMovieServerId(movie, dbo, serverAddress)
	if err != nil {
		return err
	}

	err = fetchNewMovieContainerTypeId(movie, dbo, mediaInfoResponse)
	if err != nil {
		return err
	}

	err = fetchNewMovieVideoCodecId(movie, dbo, mediaInfoResponse)
	if err != nil {
		return err
	}

	err = fetchNewMovieAudioCodecId(movie, dbo, mediaInfoResponse)
	if err != nil {
		return err
	}

	return nil
}

// fetchNewMovieServerId fetches the Server ID from the database and saves it
// into the movie object.
func fetchNewMovieServerId(
	movie *models.MovieRaw, // <- Data Receiver.
	dbo *database.DbForFilesImport,
	serverAddress string,
) (err error) {
	var serverId uint
	serverId, err = dbo.GetServerId(serverAddress)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("server is not found: %s", serverAddress)
		}
		return err
	}

	// Save the data.
	movie.FileServerId = serverId
	return nil
}

// fetchNewMovieContainerTypeId fetches the Container Type ID from the database
// and saves it into the movie object.
func fetchNewMovieContainerTypeId(
	movie *models.MovieRaw, // <- Data Receiver.
	dbo *database.DbForFilesImport,
	mediaInfoResponse *mi.Response,
) (err error) {
	var containerTypeId uint
	containerTypeId, err = dbo.GetContainerTypeId(mediaInfoResponse.Media.Tracks[0].Format)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("container is not found: %s", mediaInfoResponse.Media.Tracks[0].Format)
		}
		return err
	}

	// Save the data.
	movie.ContainerTypeId = containerTypeId
	return nil
}

// fetchNewMovieVideoCodecId fetches the Video Codec ID from the database and
// saves it into the movie object.
func fetchNewMovieVideoCodecId(
	movie *models.MovieRaw, // <- Data Receiver.
	dbo *database.DbForFilesImport,
	mediaInfoResponse *mi.Response,
) (err error) {
	var videoCodecId uint
	videoCodecId, err = dbo.GetCodecId(mediaInfoResponse.Media.Tracks[1].Format)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("video codec is not found: %s", mediaInfoResponse.Media.Tracks[1].Format)
		}
		return err
	}

	// Save the data.
	movie.VideoCodecId = videoCodecId
	return nil
}

// fetchNewMovieAudioCodecId fetches the Audio Codec ID from the database and
// saves it into the movie object.
func fetchNewMovieAudioCodecId(
	movie *models.MovieRaw, // <- Data Receiver.
	dbo *database.DbForFilesImport,
	mediaInfoResponse *mi.Response,
) (err error) {
	var audioCodecId uint
	audioCodecId, err = dbo.GetCodecId(mediaInfoResponse.Media.Tracks[2].Format)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("audio codec is not found: %s", mediaInfoResponse.Media.Tracks[2].Format)
		}
		return err
	}

	// Save the data.
	movie.AudioCodecId = audioCodecId
	return nil
}
