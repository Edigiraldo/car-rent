DROP TABLE IF EXISTS reservations;
CREATE TYPE RESERVATION_STATUSES AS ENUM('Active', 'Cancelled');
CREATE TABLE reservations (
    id uuid PRIMARY KEY NOT NULL,
    user_id uuid NOT NULL,
    car_id uuid NOT NULL REFERENCES cars(id),
    status RESERVATION_STATUSES NOT NULL DEFAULT 'Active'
);