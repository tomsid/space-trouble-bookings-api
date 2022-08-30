CREATE TABLE IF NOT EXISTS destinations (
    id serial PRIMARY KEY,
    name VARCHAR (50) UNIQUE NOT NULL
);
INSERT INTO destinations (name) VALUES
    ('Mars'),
    ('Moon'),
    ('Pluto'),
    ('Asteroid Belt'),
    ('Europa'),
    ('Titan'),
    ('Ganymede');