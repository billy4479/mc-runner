CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT,
  player_id INTEGER UNIQUE,
  tg_id INTEGER UNIQUE,
  FOREIGN KEY (player_id) REFERENCES players(id)
);

CREATE TABLE players (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  in_game_name TEXT,
  is_online BOOLEAN NOT NULL DEFAULT FALSE,
  is_bot BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE token_store (
  token BLOB PRIMARY KEY,
  expires DATE NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('auth_token', 'new_device_token', 'invitation_token', 'link_tg_token')),
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE metadata (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL
);
