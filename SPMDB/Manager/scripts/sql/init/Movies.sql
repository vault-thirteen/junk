CREATE TABLE Movies
(
    Id                 INTEGER NOT NULL PRIMARY KEY,
    Title              TEXT    NOT NULL,
    FileServerId       INTEGER NOT NULL,
    FilePath           TEXT    NOT NULL,
    FileName           TEXT    NOT NULL,
    FileExtension      TEXT    NOT NULL,
    ContainerTypeId    INTEGER NOT NULL,
    OverallBitRate     REAL    NOT NULL,
    VideoCodecId       INTEGER NOT NULL,
    AudioCodecId       INTEGER NOT NULL,
    FrameWidth         INTEGER NOT NULL,
    FrameHeight        INTEGER NOT NULL,
    Duration           INTEGER NOT NULL,
    TagsList           TEXT    NOT NULL DEFAULT '[]',
    ActorsList         TEXT    NOT NULL DEFAULT '[]',
    GenresList         TEXT    NOT NULL DEFAULT '[]',
    TimeOfRegistration INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    TimeOfUpdate       INTEGER,
    LastViewTime       INTEGER,
    ViewsCount         INTEGER NOT NULL DEFAULT 0,
    Description        TEXT    NOT NULL DEFAULT '',
    CommentsList       TEXT    NOT NULL DEFAULT '[]',
    HardcoreLevel      INTEGER NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX idx_Movies_file ON Movies (FileServerId, FilePath, FileName);
