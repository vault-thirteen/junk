CREATE TABLE IF NOT EXISTS Sessions
(
    Id        bigint AUTO_INCREMENT NOT NULL,
    UserId    bigint                NOT NULL,
    StartTime datetime              NOT NULL DEFAULT NOW(),
    UserIPAB  binary(16)            NOT NULL,
    PRIMARY KEY (Id),
    INDEX idx_UserId USING BTREE (UserId),
    INDEX idx_StartTime USING BTREE (StartTime),
    INDEX idx_UserIPAB USING BTREE (UserIPAB)
);
