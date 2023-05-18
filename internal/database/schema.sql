CREATE TABLE IF NOT EXISTS counter (
	id INTEGER PRIMARY KEY,
	count INTEGER NOT NULL
) STRICT;

INSERT OR IGNORE INTO counter (id, count) VALUES (1, 0);