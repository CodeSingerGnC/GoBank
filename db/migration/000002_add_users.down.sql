-- 假设你已经知道外键和索引的名称
ALTER TABLE `accounts` DROP FOREIGN KEY `accounts_ibfk_1`;
DROP INDEX `accounts_index_1` ON `accounts`;
DROP TABLE IF EXISTS `users`;
