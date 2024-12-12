CREATE TABLE IF NOT EXISTS UserSubscriptions
(
    Id      bigint AUTO_INCREMENT NOT NULL,
    UserId  bigint                NOT NULL,
    Threads json,

    PRIMARY KEY (Id),
    INDEX idx_UserId USING BTREE (UserId)
);
