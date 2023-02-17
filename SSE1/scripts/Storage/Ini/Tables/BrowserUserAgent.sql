CREATE TABLE IF NOT EXISTS `BrowserUserAgent`
(
	`Id` bigint unsigned NOT NULL AUTO_INCREMENT,
	`Name` text NOT NULL,
	#
	PRIMARY KEY (`Id`),
	UNIQUE KEY `Id_UNIQUE` (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
