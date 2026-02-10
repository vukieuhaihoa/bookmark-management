ALTER TABLE bookmarks
ADD COLUMN code_shorten BIGSERIAL,
ADD COLUMN code_shorten_encoded VARCHAR(25);
