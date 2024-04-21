DROP TABLE IF EXISTS readings;

DROP TABLE IF EXISTS kinds;

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

INSERT INTO
    kinds (kind_name, unit)
VALUES
    ('Gas', 'kWh'),
    ('Electricity', 'kWh');

INSERT INTO
    readings (kind_id, recorded_on, reading)
VALUES
    (1, '2023-01-01', 1234.56),
    (1, '2023-01-08', 1235.67),
    (1, '2023-01-15', 1242.67),
    (1, '2023-01-22', 1248.67),
    (1, '2023-01-29', 1252.67),
    (2, '2023-01-01', 456.8),
    (2, '2023-01-07', 526.8),
    (2, '2023-01-15', 612.8);
