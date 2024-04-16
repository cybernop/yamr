DROP DATABASE IF EXISTS readings;

CREATE DATABASE readings;

\c readings;

DROP TABLE IF EXISTS readings;

CREATE TABLE readings (
    id          SERIAL,
    kind        text NOT NULL,
    recorded_on date NOT NULL,
    reading     real NOT NULL
);

INSERT INTO readings (kind, recorded_on, reading) VALUES
    ('gas', '2023-01-01', 1234.56),
    ('gas', '2023-01-08', 1235.67),
    ('gas', '2023-01-15', 1242.67),
    ('gas', '2023-01-22', 1248.67),
    ('gas', '2023-01-29', 1252.67),
    ('electricity', '2023-01-01', 456.8),
    ('electricity', '2023-01-07', 526.8),
    ('electricity', '2023-01-15', 612.8);
