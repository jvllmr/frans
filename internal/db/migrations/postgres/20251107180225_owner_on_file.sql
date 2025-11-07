-- Modify "files" table
ALTER TABLE "files" ADD COLUMN "user_files" uuid NOT NULL, ADD CONSTRAINT "files_users_files" FOREIGN KEY ("user_files") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Drop "user_fileinfos" table
DROP TABLE "user_fileinfos";
