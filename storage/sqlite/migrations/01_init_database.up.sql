CREATE TABLE IF NOT EXISTS curse_addon (
  id INTEGER PRIMARY KEY,
  name TEXT,
  url TEXT,
  summary TEXT,
  downloadcount INTEGER
);


CREATE TABLE IF NOT EXISTS curse_release (
  id INTEGER PRIMARY KEY,
  filename TEXT NOT NULL,
  created_at INTEGER NOT NULL,
  url TEXT NOT NULL,
  addon_id INTEGER NOT NULL,
  is_alternate INTEGER NOT NULL,
  FOREIGN KEY(addon_id) REFERENCES curse_addon(id) ON DELETE CASCADE
);

CREATE INDEX curse_addon_index ON curse_release(addon_id);