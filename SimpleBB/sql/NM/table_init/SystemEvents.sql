CREATE TABLE IF NOT EXISTS SystemEvents
(
    Id        bigint AUTO_INCREMENT NOT NULL,
    Type      tinyint unsigned      NOT NULL,
    ThreadId  bigint,
    MessageId bigint,
    UserId    bigint,
    Time      datetime              NOT NULL DEFAULT NOW(),

    PRIMARY KEY (Id),
    INDEX idx_Type USING BTREE (Type),
    INDEX idx_ThreadId USING BTREE (ThreadId),
    INDEX idx_MessageId USING BTREE (MessageId),
    INDEX idx_UserId USING BTREE (UserId),
    INDEX idx_Time USING BTREE (Time)
);
