CREATE TABLE IF NOT EXISTS Users
(
    Id                bigint AUTO_INCREMENT NOT NULL,
    PreRegTime        datetime              NOT NULL,
    Email             varchar(255)          NOT NULL,
    Name              varchar(255)          NOT NULL,
    Password          varbinary(255)        NOT NULL,
    ApprovalTime      datetime              NOT NULL,
    RegTime           datetime              NOT NULL,
    IsAuthor          boolean               NOT NULL DEFAULT FALSE,
    IsWriter          boolean               NOT NULL DEFAULT FALSE,
    IsReader          boolean               NOT NULL DEFAULT FALSE,
    CanLogIn          boolean               NOT NULL DEFAULT FALSE,
    LastBadLogInTime  datetime,
    BanTime           datetime,
    LastBadActionTime datetime,
    PRIMARY KEY (Id),
    UNIQUE INDEX idx_Email USING BTREE (Email),
    UNIQUE INDEX idx_Name USING BTREE (Name),
    INDEX idx_CanLogIn USING BTREE (CanLogIn)
);
