#!/bin/bash

DB_FILE="$1"
SCHEMA_FILE="init-db.sql"

# Check if database path is provided
if [ -z "$1" ]; then
    echo "No path provided, defaulting to /var/lib/robot/metadata.sql"
    DB_FILE="/var/lib/robot/metadata.sql"
fi

DB_DIR=$(dirname "$DB_FILE")
if [ ! -d "$DB_DIR" ]; then
    echo "Directory does not exist. Creating directory: $DB_DIR"
    mkdir -p "$DB_DIR"
fi

# Check if the schema file exists
if [ ! -f "$SCHEMA_FILE" ]; then
    echo "Schema file not found: $SCHEMA_FILE"
    exit 1
fi

echo "Initializing database: $DB_FILE"
sqlite3 "$DB_FILE" < "$SCHEMA_FILE"
echo "Database initialized successfully."

