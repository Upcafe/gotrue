ALTER TABLE `refresh_tokens` ADD `updated_at` timestamp NULL DEFAULT NULL AFTER `created_at`;
UPDATE `refresh_tokens` SET `updated_at` = `created_at` WHERE `updated_at` IS NULL;
