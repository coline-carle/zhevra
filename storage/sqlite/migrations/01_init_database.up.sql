
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
  date INTEGER,
  downloadurl TEXT,
  gameversion TEXT,
  addonid INTEGER,
  FOREIGN KEY(addonid) REFERENCES curse_addon(id)
);
