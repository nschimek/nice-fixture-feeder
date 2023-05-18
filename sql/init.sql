GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, INDEX, DROP, ALTER, CREATE TEMPORARY TABLES, LOCK TABLES ON nice_fixture.* TO 'go';

CREATE TABLE nice_fixture.`leagues` (
  `id` int unsigned NOT NULL,
  `name` varchar(100) NOT NULL,
  `country_name` varchar(100) NOT NULL,
  `country_code` varchar(3) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `league_seasons` (
  `league_id` int unsigned NOT NULL,
  `season` smallint unsigned NOT NULL,
  `start` date NOT NULL,
  `end` date NOT NULL,
  `current` tinyint NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`league_id`,`season`),
  CONSTRAINT `league_seasons.league_id2leagues.id` FOREIGN KEY (`league_id`) REFERENCES `leagues` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
