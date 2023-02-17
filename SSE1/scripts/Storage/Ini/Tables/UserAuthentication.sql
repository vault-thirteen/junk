CREATE TABLE `UserAuthentication`
(
	`Id` bigint unsigned NOT NULL AUTO_INCREMENT,
	`UserId` bigint unsigned NOT NULL,
	`Name` varchar(40) NOT NULL,
	`Password` varchar(64) NOT NULL,
    `LastLogInAttemptTime` datetime DEFAULT NULL,
    #
	PRIMARY KEY (`Id`),
	UNIQUE KEY `Id_UNIQUE` (`Id`),
	UNIQUE KEY `Name_UNIQUE` (`Name`),
	UNIQUE KEY `UserId_UNIQUE` (`UserId`),
	CONSTRAINT `FK_UserAuthentication_UserId` FOREIGN KEY (`UserId`) REFERENCES `User` (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
