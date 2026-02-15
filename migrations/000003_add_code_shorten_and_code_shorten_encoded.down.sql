ALTER TABLE bookmarks
  ADD CONSTRAINT bookmarks_code_unique UNIQUE (code);

ALTER TABLE bookmarks
  ALTER COLUMN code SET NOT NULL;

  
ALTER TABLE bookmarks
  DROP COLUMN IF EXISTS code_shorten,
  DROP COLUMN IF EXISTS code_shorten_encoded;

DROP INDEX IF EXISTS idx_bookmarks_code_shorten_encoded;

ALTER TABLE bookmarks
  DROP CONSTRAINT IF EXISTS bookmarks_code_shorten_unique;
