CREATE TABLE Comments
(
    Id             INTEGER NOT NULL PRIMARY KEY,
    TimeOfCreation INTEGER NOT NULL,
    Text           TEXT    NOT NULL,
    AuthorId       INTEGER NOT NULL,
    Reputation     INTEGER NOT NULL
);
