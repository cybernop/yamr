#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE TABLE
        kinds (
            kind_id INT GENERATED ALWAYS AS IDENTITY,
            kind_name TEXT NOT NULL,
            unit TEXT NOT NULL,
            PRIMARY KEY (kind_id),
            UNIQUE (kind_name)
        );

    CREATE TABLE
        readings (
            reading_id INT GENERATED ALWAYS AS IDENTITY,
            kind_id INT NOT NULL,
            recorded_on DATE NOT NULL,
            reading REAL NOT NULL,
            PRIMARY KEY (reading_id),
            CONSTRAINT fk_kind FOREIGN KEY (kind_id) REFERENCES kinds (kind_id),
            UNIQUE (kind_id, recorded_on)
        );
EOSQL
