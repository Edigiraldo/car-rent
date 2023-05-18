DROP TABLE IF EXISTS reservations;
CREATE TYPE RESERVATION_STATUSES AS ENUM('Reserved', 'Canceled', 'Completed');
CREATE TYPE PAYMENT_STATUSES AS ENUM('Paid', 'Pending', 'Canceled');
CREATE TABLE reservations (
    id uuid PRIMARY KEY NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    car_id uuid NOT NULL REFERENCES cars(id) ON DELETE CASCADE,
    status RESERVATION_STATUSES NOT NULL DEFAULT 'Reserved',
    payment_status PAYMENT_STATUSES NOT NULL DEFAULT 'Pending',
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ
);
CREATE INDEX reservations_user_id_idx ON reservations (user_id);
CREATE INDEX reservations_car_id_idx ON reservations (car_id);
