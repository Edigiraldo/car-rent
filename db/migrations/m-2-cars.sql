DROP TABLE IF EXISTS cars;
CREATE TYPE CAR_TYPES AS ENUM('Sedan', 'Luxury', 'Sports Car', 'Limousine');
CREATE TYPE CAR_STATUSES AS ENUM('Available','Unavailable');
CREATE TABLE cars (
    id uuid NOT NULL PRIMARY KEY,
    type CAR_TYPES NOT NULL,
    seats SMALLINT NOT NULL,
    hourly_rent_cost NUMERIC(6, 2) NOT NULL,
    city_id uuid NOT NULL REFERENCES cities(id),
    status CAR_STATUSES NOT NULL
);
