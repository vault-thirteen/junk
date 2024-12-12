CREATE TABLE IF NOT EXISTS PreSessions
(
    Id                   bigint AUTO_INCREMENT NOT NULL,
    UserId               bigint                NOT NULL,
    TimeOfCreation       datetime              NOT NULL DEFAULT NOW(),
    RequestId            varchar(255)          NOT NULL,
    UserIPAB             binary(16)            NOT NULL,
    AuthDataBytes        varbinary(1024)       NOT NULL,
    IsCaptchaRequired    boolean               NOT NULL,
    CaptchaId            varchar(255), -- Null when not needed --
    IsVerifiedByCaptcha  boolean,      -- Null when not needed --
    IsVerifiedByPassword boolean               NOT NULL DEFAULT FALSE,
    VerificationCode     varchar(255)                   DEFAULT NULL,
    IsEmailSent          boolean               NOT NULL DEFAULT FALSE,
    IsVerifiedByEmail    boolean               NOT NULL DEFAULT FALSE,
    PRIMARY KEY (Id),
    INDEX idx_UserId USING BTREE (UserId),
    INDEX idx_TimeOfCreation USING BTREE (TimeOfCreation),
    INDEX idx_RequestId USING BTREE (RequestId),
    INDEX idx_IsCaptchaRequired USING BTREE (IsCaptchaRequired),
    INDEX idx_IsVerifiedByCaptcha USING BTREE (IsVerifiedByCaptcha),
    INDEX idx_IsVerifiedByPassword USING BTREE (IsVerifiedByPassword),
    INDEX idx_VerificationCode USING BTREE (VerificationCode),
    INDEX idx_IsEmailSent USING BTREE (IsEmailSent),
    INDEX idx_IsVerifiedByEmail USING BTREE (IsVerifiedByEmail)
);
