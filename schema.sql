CREATE TABLE IF NOT EXISTS accounts (
  uuid TEXT PRIMARY KEY,
  username TEXT NOT NULL,
  hashed_password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS player (
  uuid TEXT PRIMARY KEY,
  screenname TEXT,
  zone TEXT,
  x smallint,
  y smallint,
  z smallint,
  avatar smallint,
  stats json,
  items json,
  state json
);