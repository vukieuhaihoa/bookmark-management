-- The "code" column will be unused in a later migration and replaced by these two new columns
-- The "code" column is altered to allow NULL values and its unique constraint is dropped
-- The "code" column will be replaced by "code_shorten_encoded" in the application logic
ALTER TABLE bookmarks
  ALTER COLUMN code DROP NOT NULL;

ALTER TABLE bookmarks
  DROP CONSTRAINT bookmarks_code_unique;

ALTER TABLE bookmarks
  ADD COLUMN code_shorten BIGSERIAL,
  ADD COLUMN code_shorten_encoded VARCHAR(25);
