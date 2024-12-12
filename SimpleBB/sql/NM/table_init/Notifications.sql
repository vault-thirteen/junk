CREATE TABLE IF NOT EXISTS Notifications
(
    Id     bigint AUTO_INCREMENT NOT NULL,
    UserId bigint                NOT NULL,
    Text   text                  NOT NULL,

    /* ToC = Time of Creation, ToR = Time of Reading */
    ToC    datetime              NOT NULL,
    IsRead boolean               NOT NULL DEFAULT FALSE,
    ToR    datetime                       DEFAULT NULL,

    PRIMARY KEY (Id),
    INDEX idx_UserId USING BTREE (UserId),
    INDEX idx_IsRead USING BTREE (IsRead),
    INDEX idx_ToR USING BTREE (ToR)
);
