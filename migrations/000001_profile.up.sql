CREATE TABLE IF NOT EXISTS profile (
    character_id INTEGER NOT NULL UNIQUE,
    blizzard_id TEXT NOT NULL,
    battletag TEXT NOT NULL,
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

CREATE TABLE IF NOT EXISTS guild (
    character_id INTEGER NOT NULL UNIQUE REFERENCES profile(character_id),
    guild_id INTEGER NOT NULL,
    name VARCHAR(50) NOT NULL,
    name_slug VARCHAR(70) NOT NULL,
    realm VARCHAR(25) NOT NULL,
    realm_slug VARCHAR(30) NOT NULL,
    faction VARCHAR(10) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_character_blizzard_id ON profile (blizzard_id);

CREATE INDEX IF NOT EXISTS idx_character_battletag ON profile (battletag);

CREATE UNIQUE INDEX idx_profile_blizzard_id_main ON profile (blizzard_id) WHERE is_main = true;