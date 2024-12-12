CREATE TABLE IF NOT EXISTS Forums
(
    Id            bigint AUTO_INCREMENT NOT NULL,
    SectionId     bigint                NOT NULL,
    Name          varchar(255)          NOT NULL,
    Threads       json,

    -- Meta data --
    CreatorUserId bigint                NOT NULL,
    CreatorTime   datetime              NOT NULL,
    EditorUserId  bigint,
    EditorTime    datetime,

    PRIMARY KEY (Id)
    /* TODO: Indices */
);
