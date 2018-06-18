
CREATE TABLE IF NOT EXISTS curse_addon (
  id INTEGER PRIMARY KEY,
  name TEXT,
  url TEXT,
  summary TEXT,
  downloadcount INTEGER,
  addonid INTEGER,
  FOREIGN KEY(addonid) REFERENCES addon(id)
);


CREATE TABLE IF NOT EXISTS curse_release (
  id INTEGER PRIMARY KEY,
  filename TEXT,
  created_at INTEGER,
  url TEXT,
  game_version TEXT,
  addon_id INTEGER,
  FOREIGN KEY(addon_id) REFERENCES curse_addon(id)
);
