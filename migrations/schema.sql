CREATE DATABASE IF NOT EXISTS api_base;

use api_base;

CREATE TABLE IF NOT EXISTS `api_base`.`api` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
