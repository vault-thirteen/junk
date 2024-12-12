CREATE TABLE IF NOT EXISTS Resources
(
    Id     bigint AUTO_INCREMENT NOT NULL,
    Type   tinyint               NOT NULL,
    Text   text,
    Number bigint,
    ToC    datetime              NOT NULL DEFAULT NOW(),
    PRIMARY KEY (Id)
);
