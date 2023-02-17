CREATE TABLE IF NOT EXISTS `User`
(
	`Id` bigint unsigned NOT NULL AUTO_INCREMENT,
	`IsEnabled` tinyint unsigned NOT NULL DEFAULT 1,
	`PublicName` varchar(255) NOT NULL,
	#
	PRIMARY KEY (`Id`),
	UNIQUE KEY `Id_UNIQUE` (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
