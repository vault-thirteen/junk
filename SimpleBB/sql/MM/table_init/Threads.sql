CREATE TABLE IF NOT EXISTS Threads
(
    Id            bigint AUTO_INCREMENT NOT NULL,
    ForumId       bigint                NOT NULL,
    Name          varchar(255)          NOT NULL,
    Messages      json,

    -- Meta data --
    CreatorUserId bigint                NOT NULL,
    CreatorTime   datetime              NOT NULL,
    EditorUserId  bigint,
    EditorTime    datetime,

    PRIMARY KEY (Id)
    /* TODO: Indices */
);
