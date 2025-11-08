-- Modify "files" table
ALTER TABLE `files` ADD COLUMN `grant_files` char(36) NULL AFTER `file_data`, ADD COLUMN `ticket_files` char(36) NULL AFTER `grant_files`, ADD INDEX `files_grants_files` (`grant_files`), ADD INDEX `files_tickets_files` (`ticket_files`), ADD CONSTRAINT `files_grants_files` FOREIGN KEY (`grant_files`) REFERENCES `grants` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL, ADD CONSTRAINT `files_tickets_files` FOREIGN KEY (`ticket_files`) REFERENCES `tickets` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL;
-- Drop "grant_files" table
DROP TABLE `grant_files`;
-- Drop "ticket_files" table
DROP TABLE `ticket_files`;
