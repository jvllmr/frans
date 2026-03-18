-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_grants" table
CREATE TABLE `new_grants` (
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
  `email_on_upload` json NULL,
  `creator_lang` text NOT NULL DEFAULT 'en',
  `user_grants` uuid NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `grants_users_grants` FOREIGN KEY (`user_grants`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Copy rows from old table "grants" to new temporary table "new_grants"
INSERT INTO `new_grants` (`id`, `comment`, `expiry_type`, `hashed_password`, `salt`, `created_at`, `expiry_total_days`, `expiry_days_since_last_upload`, `expiry_total_uploads`, `file_expiry_type`, `file_expiry_total_days`, `file_expiry_days_since_last_download`, `file_expiry_total_downloads`, `last_upload`, `times_uploaded`, `email_on_upload`, `creator_lang`, `user_grants`) SELECT `id`, `comment`, `expiry_type`, `hashed_password`, `salt`, `created_at`, `expiry_total_days`, `expiry_days_since_last_upload`, `expiry_total_uploads`, `file_expiry_type`, `file_expiry_total_days`, `file_expiry_days_since_last_download`, `file_expiry_total_downloads`, `last_upload`, `times_uploaded`, CASE
  WHEN `email_on_upload` IS NULL THEN NULL
  ELSE json_array(`email_on_upload`)
END, `creator_lang`, `user_grants` FROM `grants`;
-- Drop "grants" table after copying rows
DROP TABLE `grants`;
-- Rename temporary table "new_grants" to "grants"
ALTER TABLE `new_grants` RENAME TO `grants`;
-- Create "new_tickets" table
CREATE TABLE `new_tickets` (
  `id` uuid NOT NULL,
  `comment` text NULL,
  `expiry_type` text NOT NULL,
  `hashed_password` text NOT NULL,
  `salt` text NOT NULL,
  `created_at` datetime NOT NULL,
  `expiry_total_days` integer NOT NULL,
  `expiry_days_since_last_download` integer NOT NULL,
  `expiry_total_downloads` integer NOT NULL,
  `email_on_download` json NULL,
  `creator_lang` text NOT NULL DEFAULT 'en',
  `user_tickets` uuid NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `tickets_users_tickets` FOREIGN KEY (`user_tickets`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Copy rows from old table "tickets" to new temporary table "new_tickets"
INSERT INTO `new_tickets` (`id`, `comment`, `expiry_type`, `hashed_password`, `salt`, `created_at`, `expiry_total_days`, `expiry_days_since_last_download`, `expiry_total_downloads`, `email_on_download`, `creator_lang`, `user_tickets`) SELECT `id`, `comment`, `expiry_type`, `hashed_password`, `salt`, `created_at`, `expiry_total_days`, `expiry_days_since_last_download`, `expiry_total_downloads`, CASE
  WHEN `email_on_download` IS NULL THEN NULL
  ELSE json_array(`email_on_download`)
END, `creator_lang`, `user_tickets` FROM `tickets`;
-- Drop "tickets" table after copying rows
DROP TABLE `tickets`;
-- Rename temporary table "new_tickets" to "tickets"
ALTER TABLE `new_tickets` RENAME TO `tickets`;
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
