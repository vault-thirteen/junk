CREATE TABLE Servers
(
    Id          INTEGER NOT NULL PRIMARY KEY,
    Address     TEXT    NOT NULL UNIQUE,
    Name        TEXT    NOT NULL UNIQUE,
    Description TEXT    NOT NULL
);
