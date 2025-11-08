-- Create "file_data" table
CREATE TABLE `file_data` (
  `id` varchar(255) NOT NULL,
  `size` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "files" table
CREATE TABLE `files` (
  `id` char(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL,
  `last_download` timestamp NULL,
  `times_downloaded` bigint unsigned NOT NULL DEFAULT 0,
  `expiry_type` varchar(255) NOT NULL,
  `expiry_total_days` tinyint unsigned NOT NULL,
  `expiry_days_since_last_download` tinyint unsigned NOT NULL,
  `expiry_total_downloads` tinyint unsigned NOT NULL,
  `file_data` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `files_file_data_data` (`file_data`),
  CONSTRAINT `files_file_data_data` FOREIGN KEY (`file_data`) REFERENCES `file_data` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "users" table
CREATE TABLE `users` (
  `id` char(36) NOT NULL,
  `username` varchar(255) NOT NULL,
  `full_name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `groups` json NOT NULL,
  `is_admin` bool NOT NULL,
  `created_at` timestamp NOT NULL,
  `submitted_tickets` bigint NOT NULL DEFAULT 0,
  `submitted_grants` bigint NOT NULL DEFAULT 0,
  `total_data_size` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "grants" table
CREATE TABLE `grants` (
  `id` char(36) NOT NULL,
  `comment` varchar(255) NULL,
  `expiry_type` varchar(255) NOT NULL,
  `hashed_password` varchar(255) NOT NULL,
  `salt` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL,
  `expiry_total_days` tinyint unsigned NOT NULL,
  `expiry_days_since_last_upload` tinyint unsigned NOT NULL,
  `expiry_total_uploads` tinyint unsigned NOT NULL,
  `file_expiry_type` varchar(255) NOT NULL,
  `file_expiry_total_days` tinyint unsigned NOT NULL,
  `file_expiry_days_since_last_download` tinyint unsigned NOT NULL,
  `file_expiry_total_downloads` tinyint unsigned NOT NULL,
  `last_upload` timestamp NULL,
  `times_uploaded` bigint unsigned NOT NULL DEFAULT 0,
  `email_on_upload` varchar(255) NULL,
  `creator_lang` varchar(255) NOT NULL DEFAULT "en",
  `user_grants` char(36) NULL,
  PRIMARY KEY (`id`),
  INDEX `grants_users_grants` (`user_grants`),
  CONSTRAINT `grants_users_grants` FOREIGN KEY (`user_grants`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "grant_files" table
CREATE TABLE `grant_files` (
  `grant_id` char(36) NOT NULL,
  `file_id` char(36) NOT NULL,
  PRIMARY KEY (`grant_id`, `file_id`),
  INDEX `grant_files_file_id` (`file_id`),
  CONSTRAINT `grant_files_file_id` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `grant_files_grant_id` FOREIGN KEY (`grant_id`) REFERENCES `grants` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "sessions" table
CREATE TABLE `sessions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `id_token` varchar(255) NOT NULL,
  `expire` timestamp NOT NULL,
  `refresh_token` varchar(255) NOT NULL,
  `user_sessions` char(36) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_token` (`id_token`),
  INDEX `sessions_users_sessions` (`user_sessions`),
  CONSTRAINT `sessions_users_sessions` FOREIGN KEY (`user_sessions`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "tickets" table
CREATE TABLE `tickets` (
  `id` char(36) NOT NULL,
  `comment` varchar(255) NULL,
  `expiry_type` varchar(255) NOT NULL,
  `hashed_password` varchar(255) NOT NULL,
  `salt` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL,
  `expiry_total_days` tinyint unsigned NOT NULL,
  `expiry_days_since_last_download` tinyint unsigned NOT NULL,
  `expiry_total_downloads` tinyint unsigned NOT NULL,
  `email_on_download` varchar(255) NULL,
  `creator_lang` varchar(255) NOT NULL DEFAULT "en",
  `user_tickets` char(36) NULL,
  PRIMARY KEY (`id`),
  INDEX `tickets_users_tickets` (`user_tickets`),
  CONSTRAINT `tickets_users_tickets` FOREIGN KEY (`user_tickets`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "share_access_tokens" table
CREATE TABLE `share_access_tokens` (
  `id` varchar(255) NOT NULL,
  `expiry` timestamp NOT NULL,
  `grant_shareaccesstokens` char(36) NULL,
  `ticket_shareaccesstokens` char(36) NULL,
  PRIMARY KEY (`id`),
  INDEX `share_access_tokens_grants_shareaccesstokens` (`grant_shareaccesstokens`),
  INDEX `share_access_tokens_tickets_shareaccesstokens` (`ticket_shareaccesstokens`),
  CONSTRAINT `share_access_tokens_grants_shareaccesstokens` FOREIGN KEY (`grant_shareaccesstokens`) REFERENCES `grants` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT `share_access_tokens_tickets_shareaccesstokens` FOREIGN KEY (`ticket_shareaccesstokens`) REFERENCES `tickets` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "ticket_files" table
CREATE TABLE `ticket_files` (
  `ticket_id` char(36) NOT NULL,
  `file_id` char(36) NOT NULL,
  PRIMARY KEY (`ticket_id`, `file_id`),
  INDEX `ticket_files_file_id` (`file_id`),
  CONSTRAINT `ticket_files_file_id` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `ticket_files_ticket_id` FOREIGN KEY (`ticket_id`) REFERENCES `tickets` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "user_fileinfos" table
CREATE TABLE `user_fileinfos` (
  `user_id` char(36) NOT NULL,
  `file_data_id` varchar(255) NOT NULL,
  PRIMARY KEY (`user_id`, `file_data_id`),
  INDEX `user_fileinfos_file_data_id` (`file_data_id`),
  CONSTRAINT `user_fileinfos_file_data_id` FOREIGN KEY (`file_data_id`) REFERENCES `file_data` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `user_fileinfos_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
