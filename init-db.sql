CREATE TABLE IF NOT EXISTS base_data (
    tag TEXT PRIMARY KEY,
    show_name TEXT,
    show_nr INTEGER,
    dj_name TEXT,
    picture TEXT,
    description TEXT,
    "tags-0-tag" TEXT,
    "tags-1-tag" TEXT,
    "tags-2-tag" TEXT,
    "tags-3-tag" TEXT,
    "tags-4-tag" TEXT
);

CREATE TABLE IF NOT EXISTS metadata (
    date TEXT,   
    tag TEXT,
    show_name TEXT,
    show_nr INTEGER,
    dj_name TEXT,
    picture TEXT,
    description TEXT,
    "tags-0-tag" TEXT,
    "tags-1-tag" TEXT,
    "tags-2-tag" TEXT,
    "tags-3-tag" TEXT,
    "tags-4-tag" TEXT,
    mixcloud BOOLEAN,
    radiocult BOOLEAN,
    drive BOOLEAN,
    local BOOLEAN,
    PRIMARY KEY (date, tag)
);

