CREATE TABLE `verify_emails` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `user_account` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `secret_code` varchar(255) NOT NULL,
  `is_used` boolean NOT NULL DEFAULT false,
  `created_at` datetime NOT NULL DEFAULT (now()),
  `expires_at` datetime NOT NULL
);

ALTER TABLE `verify_emails` ADD FOREIGN KEY (`user_account`) REFERENCES `users` (`user_account`);

ALTER TABLE `users` ADD COLUMN  `is_email_verified` boolean NOT NULL DEFAULT false;