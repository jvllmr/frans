-- Modify "grants" table
ALTER TABLE "grants" ALTER COLUMN "email_on_upload" TYPE jsonb USING CASE
  WHEN "email_on_upload" IS NULL THEN NULL
  ELSE jsonb_build_array("email_on_upload")
END;
-- Modify "tickets" table
ALTER TABLE "tickets" ALTER COLUMN "email_on_download" TYPE jsonb USING CASE
  WHEN "email_on_download" IS NULL THEN NULL
  ELSE jsonb_build_array("email_on_download")
END;
