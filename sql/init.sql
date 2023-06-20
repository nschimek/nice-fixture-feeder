GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, INDEX, DROP, ALTER, CREATE TEMPORARY TABLES, LOCK TABLES ON nice_fixture.* TO 'go';

CREATE TABLE `fixture_statuses` (
  `id` varchar(4) NOT NULL,
  `name` varchar(45) NOT NULL,
  `type` enum('SC','IP','FI','PP','CA','AB','NP') NOT NULL,
  `description` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `fixtures` (
  `id` int unsigned NOT NULL,
  `league_id` smallint unsigned NOT NULL,
  `season` smallint unsigned NOT NULL,
  `date` date NOT NULL,
  `venue_name` varchar(45) NOT NULL,
  `venue_city` varchar(45) NOT NULL,
  `status_id` varchar(4) NOT NULL,
  `team_home_id` mediumint unsigned NOT NULL,
  `team_home_result` enum('W','L','D') DEFAULT NULL,
  `team_away_id` mediumint unsigned NOT NULL,
  `team_away_result` enum('W','L','D') DEFAULT NULL,
  `goals_home` smallint unsigned DEFAULT NULL,
  `goals_away` smallint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `fixtures.league_id2tls.league_id_idx` (`league_id`,`team_home_id`,`season`),
  KEY `fixtures.team_away2team_league_seasons_idx` (`league_id`,`team_away_id`,`season`),
  KEY `fixtures.status_id2fixture_status.id_idx` (`status_id`),
  KEY `fixtures.team_home2team_league_seasons` (`team_home_id`,`league_id`,`season`),
  KEY `fixtures.team_away2team_league_seasons` (`team_away_id`,`league_id`,`season`),
  CONSTRAINT `fixtures.status_id2fixture_status.id` FOREIGN KEY (`status_id`) REFERENCES `fixture_statuses` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fixtures.team_away2team_league_seasons` FOREIGN KEY (`team_away_id`, `league_id`, `season`) REFERENCES `team_league_seasons` (`team_id`, `league_id`, `season`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fixtures.team_home2team_league_seasons` FOREIGN KEY (`team_home_id`, `league_id`, `season`) REFERENCES `team_league_seasons` (`team_id`, `league_id`, `season`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

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

CREATE TABLE `team_stats` (
  `team_id` mediumint unsigned NOT NULL,
  `league_id` smallint unsigned NOT NULL,
  `season` smallint unsigned NOT NULL,
  `fixture_id` int unsigned NOT NULL COMMENT 'the stats represent the team''s performance BEFORE playing this fixture ID; therefore, this fixture ID is not included in these stats',
  `fixtures_played_home` tinyint unsigned NOT NULL,
  `fixtures_played_away` tinyint unsigned NOT NULL,
  `fixtures_played_total` tinyint unsigned NOT NULL,
  `fixtures_wins_home` tinyint unsigned NOT NULL,
  `fixtures_wins_away` tinyint unsigned NOT NULL,
  `fixtures_wins_total` tinyint unsigned NOT NULL,
  `fixtures_draws_home` tinyint unsigned NOT NULL,
  `fixtures_draws_away` tinyint unsigned NOT NULL,
  `fixtures_draws_total` tinyint unsigned NOT NULL,
  `fixtures_losses_home` tinyint unsigned NOT NULL,
  `fixtures_losses_away` tinyint unsigned NOT NULL,
  `fixtures_losses_total` tinyint unsigned NOT NULL,
  `goals_for_home` tinyint unsigned NOT NULL,
  `goals_for_away` tinyint unsigned NOT NULL,
  `goals_for_total` tinyint unsigned NOT NULL,
  `goals_against_home` tinyint unsigned NOT NULL,
  `goals_against_away` tinyint unsigned NOT NULL,
  `goals_against_total` tinyint unsigned NOT NULL,
  `goal_differential` tinyint NOT NULL,
  `form` varchar(100) NOT NULL,
  `cs_home` tinyint unsigned NOT NULL COMMENT 'clean sheets home',
  `cs_away` tinyint unsigned NOT NULL COMMENT 'clean sheets away',
  `cs_total` tinyint unsigned NOT NULL COMMENT 'clean sheets total',
  `fts_home` tinyint unsigned NOT NULL COMMENT 'failed to score home',
  `fts_away` tinyint unsigned NOT NULL COMMENT 'failed to score away',
  `fts_total` tinyint unsigned NOT NULL COMMENT 'failed to score total',
  PRIMARY KEY (`team_id`,`league_id`,`season`,`fixture_id`),
  KEY `team_stats2team_league_seasons` (`team_id`,`league_id`,`season`),
  KEY `team_stats.fixture_id2fixtures.id_idx` (`fixture_id`),
  CONSTRAINT `team_stats.fixture_id2fixtures.id` FOREIGN KEY (`fixture_id`) REFERENCES `fixtures` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `team_stats2team_league_seasons` FOREIGN KEY (`team_id`, `league_id`, `season`) REFERENCES `team_league_seasons` (`team_id`, `league_id`, `season`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
