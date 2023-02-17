CREATE TABLE IF NOT EXISTS `Session`
(
	`Id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `UserId` bigint unsigned NOT NULL,
    `UserHost` text NOT NULL,
    `UserBuaId` bigint unsigned NOT NULL,
    `StartTime` datetime NOT NULL,
    `LastAccessTime` datetime NOT NULL DEFAULT now(),
    `EndTime` datetime DEFAULT NULL,
    `Marker` varchar(64) NOT NULL,
    `MarkerHash` varchar(64) NOT NULL,
    `TokenKey` varchar(64) NOT NULL,
    #
	PRIMARY KEY (`Id`),
	UNIQUE KEY `Id_UNIQUE` (`Id`),
    CONSTRAINT `FK_Session_UserId` FOREIGN KEY (`UserId`) REFERENCES `User` (`Id`),
    CONSTRAINT `FK_Session_UserBuaId` FOREIGN KEY (`UserBuaId`) REFERENCES `BrowserUserAgent` (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
