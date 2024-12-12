package database

import (
	"database/sql"
	"errors"
	"fmt"
	"spdb_manager/database/query"
	"spdb_manager/helper"
	"spdb_manager/models"
	"sync"

	"go.uber.org/multierr"
)

// DbSimple is a helper-object for simple operations with database.
type DbSimple struct {
	lock sync.Mutex
	db   *sql.DB
}

// Close disconnects from the database.
func (dbo *DbSimple) Close() (err error) {
	dbo.lock.Lock()
	defer dbo.lock.Unlock()

	err = dbo.db.Close()
	if err != nil {
		return err
	}

	return nil
}

// RunScript runs an SQL query from the specified script file.
func (dbo *DbSimple) RunScript(scriptFile string) (err error) {
	var q string
	q, err = helper.GetFileContentsText(scriptFile)
	if err != nil {
		return err
	}

	_, err = dbo.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

//----------------------------------------------------------------------------//

// AddActor adds an actor.
func (dbo *DbSimple) AddActor(actor *models.Actor) (actorId models.ID, err error) {
	var result sql.Result
	result, err = dbo.db.Exec(
		query.InsertActor,
		actor.Name,
		actor.Description,
	)
	if err != nil {
		return 0, err
	}

	return helper.GetSqlLastInsertId(result)
}

// AddCodec adds a codec.
func (dbo *DbSimple) AddCodec(codec *models.Codec) (codecId models.ID, err error) {
	var result sql.Result
	result, err = dbo.db.Exec(
		query.InsertCodec,
		codec.Name,
		codec.Description,
	)
	if err != nil {
		return 0, err
	}

	return helper.GetSqlLastInsertId(result)
}

// AddContainerType adds a container type.
func (dbo *DbSimple) AddContainerType(containerType *models.ContainerType) (containerTypeId models.ID, err error) {
	var result sql.Result
	result, err = dbo.db.Exec(
		query.InsertContainerTypes,
		containerType.Name,
		containerType.Description,
	)
	if err != nil {
		return 0, err
	}

	return helper.GetSqlLastInsertId(result)
}

// AddGenre adds a genre.
func (dbo *DbSimple) AddGenre(genre *models.Genre) (genreId models.ID, err error) {
	var result sql.Result
	result, err = dbo.db.Exec(
		query.InsertGenre,
		genre.Name,
		genre.Description,
	)
	if err != nil {
		return 0, err
	}

	return helper.GetSqlLastInsertId(result)
}

// AddServer adds a server.
func (dbo *DbSimple) AddServer(srv *models.Server) (serverId models.ID, err error) {
	var result sql.Result
	result, err = dbo.db.Exec(
		query.InsertServer,
		srv.Address,
		srv.Name,
		srv.Description,
	)
	if err != nil {
		return 0, err
	}

	return helper.GetSqlLastInsertId(result)
}

// AddTag adds a tag.
// It also creates an auxiliary table for the tag's movies.
func (dbo *DbSimple) AddTag(tag *models.Tag) (tagId models.ID, err error) {
	var tx *sql.Tx
	tx, err = dbo.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		var derr error
		if err != nil {
			derr = tx.Rollback()
		} else {
			derr = tx.Commit()
		}
		if derr != nil {
			err = multierr.Combine(err, derr)
		}
	}()

	var result sql.Result
	result, err = tx.Exec(
		query.InsertTag,
		tag.Name,
		tag.Description,
		tag.Colour,
	)
	if err != nil {
		return 0, err
	}

	tagId, err = helper.GetSqlLastInsertId(result)
	if err != nil {
		return 0, err
	}

	tagMoviesTableName := helper.MakeTableNameForTagNMovies(tagId)
	q := fmt.Sprintf(query.CreateTagsMoviesTableF, tagMoviesTableName)
	_, err = tx.Exec(q)
	if err != nil {
		return 0, err
	}

	return tagId, nil
}

//----------------------------------------------------------------------------//

// GetActor gets an actor.
func (dbo *DbSimple) GetActor(id models.ID) (actor *models.Actor, err error) {
	actor = new(models.Actor)
	err = dbo.db.QueryRow(query.GetActor, id).Scan(
		&actor.Id,
		&actor.Name,
		&actor.Description,
	)
	if err != nil {
		return nil, err
	}

	return actor, nil
}

// GetCodec gets a codec.
func (dbo *DbSimple) GetCodec(id models.ID) (codec *models.Codec, err error) {
	codec = new(models.Codec)
	err = dbo.db.QueryRow(query.GetCodec, id).Scan(
		&codec.Id,
		&codec.Name,
		&codec.Description,
	)
	if err != nil {
		return nil, err
	}

	return codec, nil
}

// GetContainerType gets a container type.
func (dbo *DbSimple) GetContainerType(id models.ID) (containerType *models.ContainerType, err error) {
	containerType = new(models.ContainerType)
	err = dbo.db.QueryRow(query.GetContainerType, id).Scan(
		&containerType.Id,
		&containerType.Name,
		&containerType.Description,
	)
	if err != nil {
		return nil, err
	}

	return containerType, nil
}

// GetGenre gets a genre.
func (dbo *DbSimple) GetGenre(id models.ID) (genre *models.Genre, err error) {
	genre = new(models.Genre)
	err = dbo.db.QueryRow(query.GetGenre, id).Scan(
		&genre.Id,
		&genre.Name,
		&genre.Description,
	)
	if err != nil {
		return nil, err
	}

	return genre, nil
}

// GetMovie gets a movie.
func (dbo *DbSimple) GetMovie(id models.ID) (movie *models.MovieRaw, err error) {
	movie = new(models.MovieRaw)
	err = dbo.db.QueryRow(query.GetMovie, id).Scan(
		&movie.Id,
		&movie.Title,
		&movie.FileServerId,
		&movie.FilePath,
		&movie.FileName,
		&movie.FileExtension,
		&movie.ContainerTypeId,
		&movie.OverallBitRate,
		&movie.VideoCodecId,
		&movie.AudioCodecId,
		&movie.FrameWidth,
		&movie.FrameHeight,
		&movie.Duration,
		&movie.TagsList,
		&movie.ActorsList,
		&movie.GenresList,
		&movie.TimeOfRegistration,
		&movie.TimeOfUpdate,
		&movie.LastViewTime,
		&movie.ViewsCount,
		&movie.Description,
		&movie.CommentsList,
		&movie.HardcoreLevel,
	)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

// GetServer gets a server.
func (dbo *DbSimple) GetServer(id models.ID) (server *models.Server, err error) {
	server = new(models.Server)
	err = dbo.db.QueryRow(query.GetServer, id).Scan(
		&server.Id,
		&server.Address,
		&server.Name,
		&server.Description,
	)
	if err != nil {
		return nil, err
	}

	return server, nil
}

// GetTag gets a tag.
func (dbo *DbSimple) GetTag(id models.ID) (tag *models.Tag, err error) {
	tag = new(models.Tag)
	err = dbo.db.QueryRow(query.GetTag, id).Scan(
		&tag.Id,
		&tag.Name,
		&tag.Description,
		&tag.Colour,
	)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

//----------------------------------------------------------------------------//

// SetMovieTitle changes the title of a movie specified by its ID.
func (dbo *DbSimple) SetMovieTitle(movie *models.Movie) (err error) {
	var result sql.Result
	result, err = dbo.db.Exec(
		query.SetMovieTitle,
		movie.Title,
		movie.Id,
	)
	if err != nil {
		return err
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New("no rows are affected")
	}

	return nil
}

// SetMovieDescription changes the description of a movie specified by its ID.
func (dbo *DbSimple) SetMovieDescription(movie *models.Movie) (err error) {
	var result sql.Result
	result, err = dbo.db.Exec(
		query.SetMovieDescription,
		movie.Description,
		movie.Id,
	)
	if err != nil {
		return err
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New("no rows are affected")
	}

	return nil
}

// IncMovieViewsCount increases the views count of a movie specified by its ID.
func (dbo *DbSimple) IncMovieViewsCount(movie *models.Movie) (err error) {
	var result sql.Result
	result, err = dbo.db.Exec(
		query.IncMovieViewsCount,
		movie.Id,
	)
	if err != nil {
		return err
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New("no rows are affected")
	}

	return nil
}

// SetMovieHardcoreLevel changes the hardcore level of a movie specified by its
// ID.
func (dbo *DbSimple) SetMovieHardcoreLevel(movie *models.Movie) (err error) {
	var result sql.Result
	result, err = dbo.db.Exec(
		query.SetMovieHardcoreLevel,
		movie.HardcoreLevel,
		movie.Id,
	)
	if err != nil {
		return err
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New("no rows are affected")
	}

	return nil
}

// AssignTagToMovie assigns the tag to the movie.
func (dbo *DbSimple) AssignTagToMovie(movieId models.ID, tagId models.ID) (err error) {
	var tx *sql.Tx
	tx, err = dbo.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		var derr error
		if err != nil {
			derr = tx.Rollback()
		} else {
			derr = tx.Commit()
		}
		if derr != nil {
			err = multierr.Combine(err, derr)
		}
	}()

	err = dbo.addMovieIdToTagTx(tx, movieId, tagId)
	if err != nil {
		return err
	}

	err = dbo.addTagIdToMovieTx(tx, movieId, tagId)
	if err != nil {
		return err
	}

	err = dbo.updateMovieTimeOfUpdateTx(tx, movieId)
	if err != nil {
		return err
	}

	return nil
}

// addMovieIdToTagTx adds a movie ID to the tag's table.
// This function is a part of a complex process, so it is using a transaction.
func (dbo *DbSimple) addMovieIdToTagTx(tx *sql.Tx, movieId models.ID, tagId models.ID) (err error) {
	var tagTableName = helper.MakeTableNameForTagNMovies(tagId)
	var result sql.Result
	q := fmt.Sprintf(query.AddMovieToTagF, tagTableName)
	result, err = tx.Exec(q, movieId)
	if err != nil {
		return err
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New("no rows are affected")
	}

	return nil
}

// addTagIdToMovieTx adds a tag ID to the movie's record.
// This function is a part of a complex process, so it is using a transaction.
func (dbo *DbSimple) addTagIdToMovieTx(tx *sql.Tx, movieId models.ID, tagId models.ID) (err error) {
	var movieTagsListRaw string
	err = tx.QueryRow(query.GetMovieTagsList, movieId).Scan(&movieTagsListRaw)
	if err != nil {
		return err
	}

	var movieTagsList []models.ID
	movieTagsList, err = models.ParseJsonArrayOfIds(movieTagsListRaw)
	if err != nil {
		return err
	}

	movieTagsList, err = models.AddIdToList(movieTagsList, tagId)
	if err != nil {
		return err
	}

	movieTagsListRaw, err = models.EncodeArrayOfIdsAsJson(movieTagsList)
	if err != nil {
		return err
	}

	var result sql.Result
	result, err = tx.Exec(query.SetMovieTagsList, movieTagsListRaw, movieId)
	if err != nil {
		return err
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.New("no rows are affected")
	}

	return nil
}

// updateMovieTimeOfUpdateTx updates the movie's time of update.
// This function is a part of a complex process, so it is using a transaction.
func (dbo *DbSimple) updateMovieTimeOfUpdateTx(tx *sql.Tx, movieId models.ID) (err error) {
	var result sql.Result
	result, err = tx.Exec(query.UpdateMovieTimeOfUpdate, movieId)
	if err != nil {
		return err
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.New("no rows are affected")
	}

	return nil
}

// RetractTagFromMovie retracts the tag from the movie.
func (dbo *DbSimple) RetractTagFromMovie(movieId models.ID, tagId models.ID) (err error) {
	var tx *sql.Tx
	tx, err = dbo.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		var derr error
		if err != nil {
			derr = tx.Rollback()
		} else {
			derr = tx.Commit()
		}
		if derr != nil {
			err = multierr.Combine(err, derr)
		}
	}()

	err = dbo.retractMovieIdFromTagTx(tx, movieId, tagId)
	if err != nil {
		return err
	}

	err = dbo.retractTagIdFromMovieTx(tx, movieId, tagId)
	if err != nil {
		return err
	}

	err = dbo.updateMovieTimeOfUpdateTx(tx, movieId)
	if err != nil {
		return err
	}

	return nil
}

// retractMovieIdFromTagTx retracts a movie ID to the tag's table.
// This function is a part of a complex process, so it is using a transaction.
func (dbo *DbSimple) retractMovieIdFromTagTx(tx *sql.Tx, movieId models.ID, tagId models.ID) (err error) {
	var tagTableName = helper.MakeTableNameForTagNMovies(tagId)
	var result sql.Result
	q := fmt.Sprintf(query.RemoveMovieFromTagF, tagTableName)
	result, err = tx.Exec(q, movieId)
	if err != nil {
		return err
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New("no rows are affected")
	}

	return nil
}

// retractTagIdToMovieTx retracts a tag ID from the movie's record.
// This function is a part of a complex process, so it is using a transaction.
func (dbo *DbSimple) retractTagIdFromMovieTx(tx *sql.Tx, movieId models.ID, tagId models.ID) (err error) {
	var movieTagsListRaw string
	err = tx.QueryRow(query.GetMovieTagsList, movieId).Scan(&movieTagsListRaw)
	if err != nil {
		return err
	}

	var movieTagsList []models.ID
	movieTagsList, err = models.ParseJsonArrayOfIds(movieTagsListRaw)
	if err != nil {
		return err
	}

	movieTagsList, err = models.RemoveIdFromList(movieTagsList, tagId)
	if err != nil {
		return err
	}

	movieTagsListRaw, err = models.EncodeArrayOfIdsAsJson(movieTagsList)
	if err != nil {
		return err
	}

	var result sql.Result
	result, err = tx.Exec(query.SetMovieTagsList, movieTagsListRaw, movieId)
	if err != nil {
		return err
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.New("no rows are affected")
	}

	return nil
}

// GetMoviesCount returns the count of movies.
func (dbo *DbSimple) GetMoviesCount() (count uint, err error) {
	err = dbo.db.QueryRow(query.GetMoviesCount).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
