CREATE TABLE IF NOT EXISTS Messages
(
    Id            bigint AUTO_INCREMENT NOT NULL, -- 8B --
    ThreadId      bigint                NOT NULL, -- 8B --
    Text          varchar(16368)        NOT NULL,
    TextChecksum  varbinary(4)          NOT NULL, -- 4B --

    -- Meta data --
    CreatorUserId bigint                NOT NULL, -- 8B --
    CreatorTime   datetime              NOT NULL, -- 8B --
    EditorUserId  bigint,                         -- 8B --
    EditorTime    datetime,                       -- 8B --

    PRIMARY KEY (Id)
    /* TODO: Indices */
);
