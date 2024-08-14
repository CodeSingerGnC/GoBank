-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: MySQL
-- Generated at: 2024-08-05T14:00:49.948Z

CREATE TABLE `users` (
  `user_account` varchar(255) PRIMARY KEY,
  `hash_password` varchar(255) NOT NULL,
  `username` varchar(255) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
  `is_email_verified` boolean NOT NULL DEFAULT false,
  `password_chaged_at` datetime NOT NULL DEFAULT '1970-01-01 00:00:01',
  `created_at` datetime NOT NULL DEFAULT (now())
);

CREATE TABLE `verify_emails` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `user_account` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `secret_code` varchar(255) NOT NULL,
  `is_used` boolean NOT NULL DEFAULT false,
  `created_at` datetime NOT NULL DEFAULT (now()),
  `expires_at` datetime NOT NULL DEFAULT (now() + interval 15 minutes)
);

CREATE TABLE `accounts` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `owner` varchar(255) NOT NULL,
  `balance` bigint NOT NULL DEFAULT 0,
  `currency` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (now())
);

CREATE TABLE `entries` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `account_id` bigint NOT NULL,
  `amount` bigint NOT NULL COMMENT 'can be negative or positive',
  `created_at` datetime NOT NULL DEFAULT (now())
);

CREATE TABLE `transfers` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `from_account_id` bigint NOT NULL,
  `to_account_id` bigint NOT NULL,
  `amount` bigint NOT NULL COMMENT 'must be positive',
  `created_at` datetime NOT NULL DEFAULT (now())
);

CREATE TABLE `sessions` (
  `id` BINARY(16) PRIMARY KEY,
  `user_account` varchar(255) NOT NULL,
  `refresh_token` text NOT NULL,
  `user_agent` varchar(255) NOT NULL,
  `client_ip` varchar(255) NOT NULL,
  `is_blocked` boolean NOT NULL DEFAULT false,
  `expires_at` datetime NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (now())
);

CREATE INDEX `verify_emails_index_0` ON `verify_emails` (`user_account`);

CREATE INDEX `verify_emails_index_1` ON `verify_emails` (`email`);

CREATE INDEX `accounts_index_2` ON `accounts` (`owner`);

CREATE UNIQUE INDEX `accounts_index_3` ON `accounts` (`owner`, `currency`);

CREATE INDEX `entries_index_4` ON `entries` (`account_id`);

CREATE INDEX `transfers_index_5` ON `transfers` (`from_account_id`);

CREATE INDEX `transfers_index_6` ON `transfers` (`to_account_id`);

CREATE INDEX `transfers_index_7` ON `transfers` (`from_account_id`, `to_account_id`);

ALTER TABLE `verify_emails` ADD FOREIGN KEY (`user_account`) REFERENCES `users` (`user_account`);

ALTER TABLE `accounts` ADD FOREIGN KEY (`owner`) REFERENCES `users` (`user_account`);

ALTER TABLE `entries` ADD FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`);

ALTER TABLE `transfers` ADD FOREIGN KEY (`from_account_id`) REFERENCES `accounts` (`id`);

ALTER TABLE `transfers` ADD FOREIGN KEY (`to_account_id`) REFERENCES `accounts` (`id`);

ALTER TABLE `sessions` ADD FOREIGN KEY (`user_account`) REFERENCES `users` (`user_account`);
