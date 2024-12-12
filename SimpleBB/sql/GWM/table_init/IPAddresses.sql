CREATE TABLE IF NOT EXISTS IPAddresses
(
    Id       bigint AUTO_INCREMENT NOT NULL,
    UserIPAB binary(16)                     DEFAULT NULL,
    Time     datetime              NOT NULL DEFAULT NOW(),
    PRIMARY KEY (Id),
    INDEX idx_UserIPAB USING BTREE (UserIPAB),
    INDEX idx_Time USING BTREE (Time)
);
