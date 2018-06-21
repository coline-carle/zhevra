CREATE TABLE IF NOT EXISTS curse_release_directory(
  release_id INTEGER NOT NULL,
  directory TEXT NOT NULL,
  PRIMARY KEY (release_id, directory)
  FOREIGN KEY(release_id) REFERENCES curse_release(id)
);
