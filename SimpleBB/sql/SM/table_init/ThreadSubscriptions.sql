CREATE TABLE IF NOT EXISTS ThreadSubscriptions
(
    Id       bigint AUTO_INCREMENT NOT NULL,
    ThreadId bigint                NOT NULL,
    Users    json,

    PRIMARY KEY (Id),
    INDEX idx_ThreadId USING BTREE (ThreadId)
);
