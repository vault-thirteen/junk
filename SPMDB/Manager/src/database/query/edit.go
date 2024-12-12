package query

const InsertActor = `INSERT INTO Actors (Name, Description) 
VALUES (?, ?);`

const InsertCodec = `INSERT INTO Codecs (Name, Description) 
VALUES (?, ?);`

const InsertContainerTypes = `INSERT INTO ContainerTypes (Name, Description) 
VALUES (?, ?);`

const InsertGenre = `INSERT INTO Genres (Name, Description) 
VALUES (?, ?);`

const InsertServer = `INSERT INTO Servers (Address, Name, Description) 
VALUES (?, ?, ?);`

const InsertTag = `INSERT INTO Tags (Name, Description, Colour) 
VALUES (?, ?, ?);`

const InsertNewMovie = `INSERT INTO Movies (Title, FileServerId, FilePath, ` +
	`FileName, FileExtension, ContainerTypeId, OverallBitRate, VideoCodecId, ` +
	`AudioCodecId, FrameWidth, FrameHeight, Duration) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"`

const CreateTagsMoviesTableF = `CREATE TABLE %s 
(
	MovieId INTEGER NOT NULL PRIMARY KEY
);`

const SetMovieTitle = `UPDATE Movies 
SET Title = ?, TimeOfUpdate = STRFTIME('%s', 'now') 
WHERE Id = ?;`

const SetMovieDescription = `UPDATE Movies 
SET Description = ?, TimeOfUpdate = STRFTIME('%s', 'now') 
WHERE Id = ?;`

const IncMovieViewsCount = `UPDATE Movies 
SET ViewsCount = ViewsCount + 1, LastViewTime = STRFTIME('%s', 'now') 
WHERE Id = ?;`

const SetMovieHardcoreLevel = `UPDATE Movies 
SET HardcoreLevel = ?, TimeOfUpdate = STRFTIME('%s', 'now') 
WHERE Id = ?;`

const AddMovieToTagF = `INSERT INTO %s (MovieId) VALUES (?);`

const SetMovieTagsList = `UPDATE Movies 
SET TagsList = ? 
WHERE Id = ?;`

const UpdateMovieTimeOfUpdate = `UPDATE Movies 
SET TimeOfUpdate = STRFTIME('%s', 'now') 
WHERE Id = ?;`

const RemoveMovieFromTagF = `DELETE FROM %s WHERE MovieId = ?;`
