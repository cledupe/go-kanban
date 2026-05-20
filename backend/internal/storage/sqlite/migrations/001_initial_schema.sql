CREATE TABLE IF NOT EXISTS boards (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS columns (
    id TEXT PRIMARY KEY,
    board_id TEXT NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    position INTEGER NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS cards (
    id TEXT PRIMARY KEY,
    column_id TEXT NOT NULL REFERENCES columns(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    position INTEGER NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
