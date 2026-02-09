ALTER TABLE bookmarks
  DROP COLUMN IF EXISTS code_shorten,
  DROP COLUMN IF EXISTS code_shorten_encoded;