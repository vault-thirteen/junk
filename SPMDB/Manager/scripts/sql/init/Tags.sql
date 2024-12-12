CREATE TABLE Tags
(
    Id          INTEGER NOT NULL PRIMARY KEY,
    Name        TEXT    NOT NULL UNIQUE,
    Description INTEGER NOT NULL,
    Colour      TEXT    NOT NULL
);
