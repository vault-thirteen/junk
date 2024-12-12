CREATE TABLE IF NOT EXISTS Incidents
(
    Id       bigint AUTO_INCREMENT NOT NULL,
    Module   tinyint               NOT NULL,
    Type     tinyint               NOT NULL,
    Time     datetime              NOT NULL DEFAULT NOW(),
    Email    varchar(255)          NOT NULL,
    UserIPAB binary(16)                     DEFAULT NULL,
    PRIMARY KEY (Id),
    INDEX idx_Module USING BTREE (Module),
    INDEX idx_Type USING BTREE (Type),
    INDEX idx_Time USING BTREE (Time),
    INDEX idx_Email USING BTREE (Email),
    INDEX idx_UserIPAB USING BTREE (UserIPAB)
);
