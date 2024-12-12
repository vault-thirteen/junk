CREATE TABLE IF NOT EXISTS Sections
(
    Id            bigint AUTO_INCREMENT NOT NULL,
    Parent        bigint,
    ChildType     tinyint DEFAULT 3,
    Children      json,
    Name          varchar(255)          NOT NULL,

    -- Meta data --
    CreatorUserId bigint                NOT NULL,
    CreatorTime   datetime              NOT NULL,
    EditorUserId  bigint,
    EditorTime    datetime,

    PRIMARY KEY (Id),
    INDEX idx_Parent USING BTREE (Parent)
    /* TODO: Indices */
);
