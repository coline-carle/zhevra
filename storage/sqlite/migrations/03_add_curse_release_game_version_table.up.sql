CREATE TABLE IF NOT EXISTS curse_release_game_version(
  release_id INTEGER NOT NULL,
  game_version INTEGER NOT NULL,
  PRIMARY KEY (release_id, game_version)
  FOREIGN KEY(release_id) REFERENCES curse_release(id) ON DELETE CASCADE
);

CREATE INDEX curse_release_index_gv ON curse_release_directory(release_id);