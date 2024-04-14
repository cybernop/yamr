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
    ('gas', '2023-01-02', 1235.67),
    ('electricity', '2023-01-01', 456.8);
