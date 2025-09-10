CREATE TABLE IF NOT EXISTS profile (
    blizzard_id TEXT,
    battletag TEXT,
    name VARCHAR(25) NOT NULL,
    realm VARCHAR(25) NOT NULL,
    race VARCHAR(25) NOT NULL,
    faction VARCHAR(10) NOT NULL,
    class VARCHAR(30) NOT NULL,
    spec VARCHAR(25),
    lvl INTEGER NOT NULL,
    ilvl INTEGER,
    guild TEXT,
    mythic_score NUMERIC(10,2),
    is_main BOOLEAN,
    PRIMARY KEY (blizzard_id, name, realm)
);

CREATE INDEX IF NOT EXISTS idx_character_blizzard_id ON profile (blizzard_id);

CREATE INDEX IF NOT EXISTS idx_character_battletag ON profile (battletag);

CREATE UNIQUE INDEX idx_profile_blizzard_id_main ON profile (blizzard_id) WHERE is_main = true;