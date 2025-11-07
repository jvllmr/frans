-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_files" table
CREATE TABLE `new_files` (
  `id` uuid NOT NULL,
  `name` text NOT NULL,
  `created_at` datetime NOT NULL,
  `last_download` datetime NULL,
  `times_downloaded` integer NOT NULL DEFAULT 0,
  `expiry_type` text NOT NULL,
  `expiry_total_days` integer NOT NULL,
  `expiry_days_since_last_download` integer NOT NULL,
  `expiry_total_downloads` integer NOT NULL,
  `file_data` text NOT NULL,
  `user_files` uuid NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `files_users_files` FOREIGN KEY (`user_files`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `files_file_data_data` FOREIGN KEY (`file_data`) REFERENCES `file_data` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Copy rows from old table "files" to new temporary table "new_files"
INSERT INTO `new_files` (`id`, `name`, `created_at`, `last_download`, `times_downloaded`, `expiry_type`, `expiry_total_days`, `expiry_days_since_last_download`, `expiry_total_downloads`, `file_data`) SELECT `id`, `name`, `created_at`, `last_download`, `times_downloaded`, `expiry_type`, `expiry_total_days`, `expiry_days_since_last_download`, `expiry_total_downloads`, `file_data` FROM `files`;
-- Drop "files" table after copying rows
DROP TABLE `files`;
-- Rename temporary table "new_files" to "files"
ALTER TABLE `new_files` RENAME TO `files`;
-- Drop "user_fileinfos" table
DROP TABLE `user_fileinfos`;
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
