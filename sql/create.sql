CREATE DATABASE IF NOT EXISTS crescent_roulette_db;
use crescent_roulette_db;

CREATE TABLE IF NOT EXISTS `account` (
  `address` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `voucher` int NOT NULL DEFAULT '0',
  `ticket` int NOT NULL DEFAULT '0',
  `hex_addr` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_ts` bigint NOT NULL DEFAULT '0',
  `admin_memo` tinytext COLLATE utf8mb4_unicode_ci,
  `admin_alias` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `update_ts` bigint DEFAULT NULL,
  PRIMARY KEY (`address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `game_order` (
  `game_order_id` bigint NOT NULL AUTO_INCREMENT,
  `address` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `paid_ticket_num` int NOT NULL DEFAULT '1',
  `type` int NOT NULL DEFAULT '0',
  `status` int NOT NULL DEFAULT '1',
  `is_win` boolean NOT NULL default false,
  `prize_id` int,
  `created_ts` bigint NOT NULL DEFAULT '0',
  `update_ts` bigint DEFAULT NULL,
  PRIMARY KEY (`game_order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `prize` (
  `prize_id` bigint NOT NULL AUTO_INCREMENT,
  `type` int NOT NULL DEFAULT '0',
  `token_denom` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `token_num` int NOT NULL DEFAULT '0',
  `created_ts` bigint NOT NULL DEFAULT '0',
  `update_ts` bigint DEFAULT NULL,
  PRIMARY KEY (`prize_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
