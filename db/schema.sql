CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  player_id INTEGER UNIQUE,
  tg_id INTEGER UNIQUE,
  FOREIGN KEY (player_id) REFERENCES players(id)
);

CREATE TABLE players (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  in_game_name TEXT NOT NULL,
  is_online BOOLEAN NOT NULL,
  is_bot BOOLEAN NOT NULL
);

CREATE TABLE token_store (
  token BLOB PRIMARY KEY,
  expires DATE NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('refresh_token', 'login_token', 'register_token', 'link_tg_token')),
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE metadata (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL
);
