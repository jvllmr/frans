-- Create "files" table
CREATE TABLE `files` (
  `id` uuid NOT NULL,
  `name` text NOT NULL,
  `size` integer NOT NULL,
  `sha512` text NOT NULL,
  `created_at` datetime NOT NULL,
  `last_download` datetime NULL,
  `times_downloaded` integer NOT NULL DEFAULT 0,
  `expiry_type` text NOT NULL,
  `expiry_total_days` integer NOT NULL,
  `expiry_days_since_last_download` integer NOT NULL,
  `expiry_total_downloads` integer NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create "grants" table
CREATE TABLE `grants` (
  `id` uuid NOT NULL,
  `comment` text NULL,
  `expiry_type` text NOT NULL,
  `hashed_password` text NOT NULL,
  `salt` text NOT NULL,
  `created_at` datetime NOT NULL,
  `expiry_total_days` integer NOT NULL,
  `expiry_days_since_last_upload` integer NOT NULL,
  `expiry_total_uploads` integer NOT NULL,
  `file_expiry_type` text NOT NULL,
  `file_expiry_total_days` integer NOT NULL,
  `file_expiry_days_since_last_download` integer NOT NULL,
  `file_expiry_total_downloads` integer NOT NULL,
  `last_upload` datetime NULL,
  `times_uploaded` integer NOT NULL DEFAULT 0,
  `email_on_upload` text NULL,
  `creator_lang` text NOT NULL DEFAULT 'en',
  `user_grants` uuid NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `grants_users_grants` FOREIGN KEY (`user_grants`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create "sessions" table
CREATE TABLE `sessions` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `id_token` text NOT NULL,
  `expire` datetime NOT NULL,
  `refresh_token` text NOT NULL,
  `user_sessions` uuid NULL,
  CONSTRAINT `sessions_users_sessions` FOREIGN KEY (`user_sessions`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "sessions_id_token_key" to table: "sessions"
CREATE UNIQUE INDEX `sessions_id_token_key` ON `sessions` (`id_token`);
-- Create "share_access_tokens" table
CREATE TABLE `share_access_tokens` (
  `id` text NOT NULL,
  `expiry` datetime NOT NULL,
  `grant_shareaccesstokens` uuid NULL,
  `ticket_shareaccesstokens` uuid NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `share_access_tokens_tickets_shareaccesstokens` FOREIGN KEY (`ticket_shareaccesstokens`) REFERENCES `tickets` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT `share_access_tokens_grants_shareaccesstokens` FOREIGN KEY (`grant_shareaccesstokens`) REFERENCES `grants` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create "tickets" table
CREATE TABLE `tickets` (
  `id` uuid NOT NULL,
  `comment` text NULL,
  `expiry_type` text NOT NULL,
  `hashed_password` text NOT NULL,
  `salt` text NOT NULL,
  `created_at` datetime NOT NULL,
  `expiry_total_days` integer NOT NULL,
  `expiry_days_since_last_download` integer NOT NULL,
  `expiry_total_downloads` integer NOT NULL,
  `email_on_download` text NULL,
  `creator_lang` text NOT NULL DEFAULT 'en',
  `user_tickets` uuid NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `tickets_users_tickets` FOREIGN KEY (`user_tickets`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create "users" table
CREATE TABLE `users` (
  `id` uuid NOT NULL,
  `username` text NOT NULL,
  `full_name` text NOT NULL,
  `email` text NOT NULL,
  `groups` json NOT NULL,
  `is_admin` bool NOT NULL,
  `created_at` datetime NOT NULL,
  `submitted_tickets` integer NOT NULL DEFAULT 0,
  `submitted_grants` integer NOT NULL DEFAULT 0,
  `total_data_size` integer NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
);
-- Create "grant_files" table
CREATE TABLE `grant_files` (
  `grant_id` uuid NOT NULL,
  `file_id` uuid NOT NULL,
  PRIMARY KEY (`grant_id`, `file_id`),
  CONSTRAINT `grant_files_file_id` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `grant_files_grant_id` FOREIGN KEY (`grant_id`) REFERENCES `grants` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "ticket_files" table
CREATE TABLE `ticket_files` (
  `ticket_id` uuid NOT NULL,
  `file_id` uuid NOT NULL,
  PRIMARY KEY (`ticket_id`, `file_id`),
  CONSTRAINT `ticket_files_file_id` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `ticket_files_ticket_id` FOREIGN KEY (`ticket_id`) REFERENCES `tickets` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);
