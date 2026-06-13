CREATE TABLE IF NOT EXISTS pets (
    id TEXT PRIMARY KEY,
    species TEXT NOT NULL DEFAULT 'cat',
    name TEXT NOT NULL,
    breed TEXT NOT NULL DEFAULT '',
    gender TEXT NOT NULL DEFAULT '',
    birthday TEXT NOT NULL,
    color TEXT NOT NULL DEFAULT '',
    avatar_url TEXT NOT NULL DEFAULT '',
    dominant_color TEXT NOT NULL DEFAULT '',
    adopted_at TEXT,
    passed_at TEXT,
    note TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS weight_records (
    id TEXT PRIMARY KEY,
    pet_id TEXT NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    weight REAL NOT NULL,
    recorded_at TEXT NOT NULL,
    note TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS health_records (
    id TEXT PRIMARY KEY,
    pet_id TEXT NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    name TEXT NOT NULL,
    date TEXT NOT NULL,
    next_date TEXT,
    note TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS photos (
    id TEXT PRIMARY KEY,
    pet_id TEXT NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    age_group TEXT NOT NULL,
    caption TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_weight_records_pet_id ON weight_records(pet_id);
CREATE INDEX IF NOT EXISTS idx_weight_records_recorded_at ON weight_records(pet_id, recorded_at);
CREATE INDEX IF NOT EXISTS idx_health_records_pet_id ON health_records(pet_id);
CREATE INDEX IF NOT EXISTS idx_health_records_type ON health_records(pet_id, type);
CREATE INDEX IF NOT EXISTS idx_health_records_next_date ON health_records(next_date);
CREATE INDEX IF NOT EXISTS idx_photos_pet_id ON photos(pet_id);
CREATE INDEX IF NOT EXISTS idx_photos_age_group ON photos(pet_id, age_group);
