CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT UNIQUE
);

CREATE TABLE token_store (
  token BLOB PRIMARY KEY,
  expires DATE NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('auth_token', 'login_token', 'invitation_token')),
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
