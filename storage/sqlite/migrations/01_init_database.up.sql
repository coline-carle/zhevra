
CREATE TABLE IF NOT EXISTS curse_addon (
  id INTEGER PRIMARY KEY,
  name TEXT,
  url TEXT,
  summary TEXT,
  downloadcount INTEGER
);


CREATE TABLE IF NOT EXISTS curse_release (
  id INTEGER NOT NULL,
  filename TEXT,
  created_at INTEGER,
  url TEXT NOT NULL,
  game_version TEXT,
  addon_id INTEGER NOT NULL,
  PRIMARY KEY (id, addon_id),
  FOREIGN KEY(addon_id) REFERENCES curse_addon(id)
);
