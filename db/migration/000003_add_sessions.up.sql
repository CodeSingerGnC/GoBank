CREATE TABLE `sessions` (
  `id` BINARY(16) PRIMARY KEY,
  `user_account` varchar(255) NOT NULL,
  `refresh_token` TEXT NOT NULL,
  `user_agent` varchar(255) NOT NULL,
  `client_ip` varchar(255) NOT NULL,
  `is_blocked` boolean NOT NULL DEFAULT false,
  `expires_at` datetime NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (now())
);

ALTER TABLE `sessions` ADD FOREIGN KEY (`user_account`) REFERENCES `users` (`user_account`);