CREATE TABLE `UserRegistration`
(
	`Id` bigint unsigned NOT NULL AUTO_INCREMENT,
	`UserId` bigint unsigned NOT NULL,
	`SecretCode` varchar(64) NOT NULL,
	`RegTime` datetime NOT NULL DEFAULT now(),
	`UnregTime` datetime DEFAULT NULL,
    `LastUnregAttemptTime` datetime DEFAULT NULL,
    #
	PRIMARY KEY (`Id`),
	UNIQUE KEY `Id_UNIQUE` (`Id`),
	UNIQUE KEY `UserId_UNIQUE` (`UserId`),
	CONSTRAINT `FK_UserRegistration_UserId` FOREIGN KEY (`UserId`) REFERENCES `User` (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
