CREATE TABLE ContainerTypes
(
    Id          INTEGER NOT NULL PRIMARY KEY,
    Name        TEXT    NOT NULL UNIQUE,
    Description TEXT    NOT NULL
);
