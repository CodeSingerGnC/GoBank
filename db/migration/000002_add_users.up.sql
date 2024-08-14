CREATE TABLE `users` (
  `user_account` varchar(255) PRIMARY KEY,
  `hash_password` varchar(255) NOT NULL,
  `username` varchar(255) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
  `password_chaged_at` datetime NOT NULL DEFAULT '1970-01-01 00:00:01',
  `created_at` datetime NOT NULL DEFAULT (now())
);

ALTER TABLE `accounts` ADD FOREIGN KEY (`owner`) REFERENCES `users` (`user_account`);

CREATE UNIQUE INDEX `accounts_index_1` ON `accounts` (`owner`, `currency`);