CREATE TABLE IF NOT EXISTS base_data (
    tag TEXT PRIMARY KEY,
    show_name TEXT,
    show_nr INTEGER,
    dj_name TEXT,
    picture TEXT,
    description TEXT,
    tags0 TEXT,
    tags1 TEXT,
    tags2 TEXT,
    tags3 TEXT,
    tags4 TEXT,
    live BOOLEAN
);

CREATE TABLE IF NOT EXISTS metadata (
    date TEXT,   
    tag TEXT,
    show_name TEXT,
    show_nr INTEGER,
    dj_name TEXT,
    picture TEXT,
    description TEXT,
    tags0 TEXT,
    tags1 TEXT,
    tags2 TEXT,
    tags3 TEXT,
    tags4 TEXT,
    live BOOLEAN,
    mixcloud BOOLEAN,
    soundcloud BOOLEAN,
    soundcloud_urn TEXT,
    radiocult BOOLEAN,
    drive BOOLEAN,
    PRIMARY KEY (date, tag)
);
