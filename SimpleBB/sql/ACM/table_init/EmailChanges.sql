CREATE TABLE IF NOT EXISTS EmailChanges
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
    VerificationCodeOld  varchar(255)                   DEFAULT NULL,
    IsOldEmailSent       boolean               NOT NULL DEFAULT FALSE,
    IsVerifiedByOldEmail boolean               NOT NULL DEFAULT FALSE,
    NewEmail             varbinary(255)        NOT NULL,
    VerificationCodeNew  varchar(255)                   DEFAULT NULL,
    IsNewEmailSent       boolean               NOT NULL DEFAULT FALSE,
    IsVerifiedByNewEmail boolean               NOT NULL DEFAULT FALSE,
    PRIMARY KEY (Id),
    INDEX idx_UserId USING BTREE (UserId),
    INDEX idx_TimeOfCreation USING BTREE (TimeOfCreation),
    INDEX idx_RequestId USING BTREE (RequestId),
    INDEX idx_IsCaptchaRequired USING BTREE (IsCaptchaRequired),
    INDEX idx_IsVerifiedByCaptcha USING BTREE (IsVerifiedByCaptcha),
    INDEX idx_IsVerifiedByPassword USING BTREE (IsVerifiedByPassword),
    INDEX idx_VerificationCodeOld USING BTREE (VerificationCodeOld),
    INDEX idx_VerificationCodeNew USING BTREE (VerificationCodeNew),
    INDEX idx_IsOldEmailSent USING BTREE (IsOldEmailSent),
    INDEX idx_IsNewEmailSent USING BTREE (IsNewEmailSent),
    INDEX idx_IsVerifiedByOldEmail USING BTREE (IsVerifiedByOldEmail),
    INDEX idx_IsVerifiedByNewEmail USING BTREE (IsVerifiedByNewEmail)
);
