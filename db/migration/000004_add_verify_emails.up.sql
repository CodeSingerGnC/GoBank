CREATE TABLE `otpsecrets` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `secret` varchar(64) NOT NULL,
  `tried_times` int NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT (now())
);

ALTER TABLE `otpsecrets` ADD INDEX `idx_email` (`email`);