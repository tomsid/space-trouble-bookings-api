CREATE TABLE IF NOT EXISTS bookings (
                                        id serial PRIMARY KEY,
                                        first_name VARCHAR (50) UNIQUE NOT NULL,
    last_name VARCHAR (50) NOT NULL,
    gender VARCHAR (10) NOT NULL,
    birthday date NOT NULL,
    launchpadID int NOT NULL,
    destinationID int NOT NULL,
    launchDate date NOT NULL
    );