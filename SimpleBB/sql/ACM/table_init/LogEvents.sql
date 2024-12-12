CREATE TABLE IF NOT EXISTS LogEvents
(
    Id       bigint AUTO_INCREMENT NOT NULL,
    Time     datetime              NOT NULL DEFAULT NOW(),
    Type     tinyint               NOT NULL,
    UserId   bigint                NOT NULL,
    Email    varchar(255)          NOT NULL,
    UserIPAB binary(16)                     DEFAULT NULL,
    AdminId  bigint                         DEFAULT NULL,
    PRIMARY KEY (Id),
    INDEX idx_Time USING BTREE (Time),
    INDEX idx_Type USING BTREE (Type),
    INDEX idx_UserId USING BTREE (UserId),
    INDEX idx_Email USING BTREE (Email),
    INDEX idx_UserIPAB USING BTREE (UserIPAB),
    INDEX idx_AdminId USING BTREE (AdminId)
);
