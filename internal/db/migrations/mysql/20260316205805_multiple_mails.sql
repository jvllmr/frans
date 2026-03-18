-- Modify "grants" table
UPDATE `grants`
SET `email_on_upload` =
  IF(`email_on_upload` IS NULL, NULL, JSON_QUOTE(`email_on_upload`));
ALTER TABLE `grants` MODIFY COLUMN `email_on_upload` json NULL;
UPDATE `grants`
SET `email_on_upload` =
  IF(`email_on_upload` IS NULL, NULL, JSON_ARRAY(JSON_UNQUOTE(`email_on_upload`)));
-- Modify "tickets" table
UPDATE `tickets`
SET `email_on_download` =
  IF(`email_on_download` IS NULL, NULL, JSON_QUOTE(`email_on_download`));
ALTER TABLE `tickets` MODIFY COLUMN `email_on_download` json NULL;
UPDATE `tickets`
SET `email_on_download` =
  IF(`email_on_download` IS NULL, NULL, JSON_ARRAY(JSON_UNQUOTE(`email_on_download`)));
