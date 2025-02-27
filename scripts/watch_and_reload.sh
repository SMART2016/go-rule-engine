#!/bin/bash

# Configuration
SQL_FILE="$(pwd)/store/sqlc/store_schema.sql"   # Path to the SQL file
DB_CONTAINER="postgres"                         # Name of the running PostgreSQL container
DB_USER="dbuser"                                # PostgreSQL user
DB_NAME="rule_engine"                           # Database name

echo $SQL_FILE
# Detect OS (Linux or Mac)
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    WATCH_CMD="inotifywait -e modify -q $SQL_FILE"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    WATCH_CMD="fswatch -o $SQL_FILE"
else
    echo "Unsupported OS: $OSTYPE"
    exit 1
fi

echo "Watching $SQL_FILE for changes..."
while true; do
    $WATCH_CMD
    echo "Detected change in $SQL_FILE. Reloading into PostgreSQL..."
    docker exec -i $DB_CONTAINER psql -U $DB_USER -d $DB_NAME < "$SQL_FILE"
    echo "Schema reloaded successfully."
done
