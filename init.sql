DROP TABLE IF EXISTS `houses`;
CREATE TABLE `houses` (
  `id` char(8) COLLATE utf8mb4_bin NOT NULL,
  `view_count` INT(6) NOT NULL,
  `html` text COLLATE utf8mb4_bin NOT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
