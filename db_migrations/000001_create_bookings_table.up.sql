CREATE TABLE IF NOT EXISTS bookings (
    id serial PRIMARY KEY,
    first_name VARCHAR (50) NOT NULL,
    last_name VARCHAR (50) NOT NULL,
    gender VARCHAR (10) NOT NULL,
    birthday date NOT NULL,
    launchpad_id varchar (30) NOT NULL,
    destination_id int NOT NULL,
    launch_date date NOT NULL
);