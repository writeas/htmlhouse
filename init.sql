DROP TABLE IF EXISTS `houses`;
CREATE TABLE `houses` (
  `id` char(8) COLLATE utf8mb4_bin NOT NULL,
  `view_count` INT(6) NOT NULL,
  `html` text COLLATE utf8mb4_bin NOT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

DROP TABLE IF EXISTS `publichouses`;
CREATE TABLE `publichouses` (
  `house_id` char(8) NOT NULL,
  `title` varchar(128) COLLATE utf8mb4_bin NOT NULL,
  `thumb_url` varchar(28) NOT NULL,
  `added` datetime NOT NULL,
  `updated` datetime NOT NULL,
  `approved` tinyint(1) DEFAULT NULL,
  `loves` int(4) NOT NULL,
  PRIMARY KEY (`house_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `tweetedhouses` (
  `house_id` char(8) NOT NULL,
  `tweet_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`house_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
