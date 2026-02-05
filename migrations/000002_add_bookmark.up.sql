CREATE TABLE bookmarks 
(
  id varchar(36),
  description varchar(256),
  url varchar(2048) NOT NULL,
  code varchar(10) NOT NULL,
  user_id varchar(36) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,

  CONSTRAINT bookmarks_pk PRIMARY KEY (id),
  CONSTRAINT bookmarks_code_unique UNIQUE (code),
  CONSTRAINT bookmarks_user_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);