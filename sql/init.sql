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

CREATE TABLE `teams` (
  `id` mediumint unsigned NOT NULL,
  `name` varchar(100) NOT NULL,
  `code` varchar(3) NOT NULL,
  `country` varchar(100) NOT NULL,
  `national` tinyint unsigned NOT NULL,
  `venue_name` varchar(100) NOT NULL,
  `venue_city` varchar(100) NOT NULL,
  `venue_capacity` mediumint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `team_id_UNIQUE` (`id`),
  UNIQUE KEY `code_UNIQUE` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `team_league_seasons` (
  `team_id` mediumint unsigned NOT NULL,
  `league_id` smallint unsigned NOT NULL,
  `season` smallint unsigned NOT NULL,
  PRIMARY KEY (`team_id`,`league_id`,`season`),
  KEY `team_league_seasons2league_seasons_idx` (`league_id`,`season`),
  CONSTRAINT `team_league_seasons.team_id2teams.id` FOREIGN KEY (`team_id`) REFERENCES `teams` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `team_league_seasons2league_seasons` FOREIGN KEY (`league_id`, `season`) REFERENCES `league_seasons` (`league_id`, `season`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
