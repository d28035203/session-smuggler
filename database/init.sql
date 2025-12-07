-- Optional bootstrap schema. GORM AutoMigrate also creates tables at runtime.
CREATE TABLE IF NOT EXISTS users (
    username TEXT PRIMARY KEY,
    password BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS usersessions (
    username TEXT PRIMARY KEY REFERENCES users (username) ON DELETE CASCADE,
    token TEXT NOT NULL
);
