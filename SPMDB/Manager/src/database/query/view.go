package query

const GetActor = `SELECT * FROM Actors WHERE Id = ?;`

const GetCodec = `SELECT * FROM Codecs WHERE Id = ?;`

const GetContainerType = `SELECT * FROM ContainerTypes WHERE Id = ?;`

const GetGenre = `SELECT * FROM Genres WHERE Id = ?;`

const GetMovie = `SELECT * FROM Movies WHERE Id = ?;`

const GetServer = `SELECT * FROM Servers WHERE Id = ?;`

const GetTag = `SELECT * FROM Tags WHERE Id = ?;`

const GetServerIdByAddress = `SELECT Id FROM Servers WHERE Address = ?;`

const GetContainerTypeIdByName = `SELECT Id FROM ContainerTypes WHERE Name = ?;`

const GetCodecIdByName = `SELECT Id FROM Codecs WHERE Name = ?;`

const GetMovieTagsList = `SELECT TagsList FROM Movies WHERE Id = ?;`

const GetMoviesCount = `SELECT count(Id) FROM Movies;`
