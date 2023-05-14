-- Enable the pgcrypto extension to generate random UUIDs
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create a temporary table to store the generated reservations
CREATE TEMP TABLE temp_reservations (
    id uuid PRIMARY KEY,
    user_id uuid,
    car_id uuid,
    status RESERVATION_STATUSES,
    payment_status PAYMENT_STATUSES,
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ
);

-- Generate 1000 random reservations
DO $$
DECLARE
    temp_user_id UUID;
    temp_car_id UUID;
    temp_start_date TIMESTAMPTZ;
    temp_end_date TIMESTAMPTZ;
BEGIN
    FOR i IN 1..1000 LOOP
        -- Select a random user and car
        SELECT users.id, cars.id
        INTO temp_user_id, temp_car_id
        FROM users, cars
        ORDER BY RANDOM()
        LIMIT 1;

        -- Calculate start_date and end_date
        temp_start_date := now() + INTERVAL '1 day' * floor(random() * 182); -- 6 months
        temp_end_date := temp_start_date + INTERVAL '1 day' * floor(random() * 10) + INTERVAL '7 hours';

        -- Ensure there are no overlapping reservations
        WHILE EXISTS (
            SELECT 1
            FROM reservations
            WHERE (reservations.user_id = temp_user_id OR reservations.car_id = temp_car_id)
                AND (reservations.start_date, reservations.end_date) OVERLAPS (temp_start_date, temp_end_date)
        ) LOOP
            temp_start_date := now() + INTERVAL '1 day' * floor(random() * 182); -- 6 months
            temp_end_date := temp_start_date + INTERVAL '1 day' * floor(random() * 10) + INTERVAL '7 hours';
        END LOOP;

        -- Insert the reservation into the temporary table
        INSERT INTO temp_reservations (id, user_id, car_id, status, payment_status, start_date, end_date)
        VALUES (
            gen_random_uuid(),
            temp_user_id,
            temp_car_id,
            'Reserved',
            'Pending',
            temp_start_date,
            temp_end_date
        );
    END LOOP;
END;
$$;

-- Insert the generated reservations into the main table
INSERT INTO reservations (id, user_id, car_id, status, payment_status, start_date, end_date)
SELECT id, user_id, car_id, status, payment_status, start_date, end_date
FROM temp_reservations;

-- Clean up the temporary table
DROP TABLE temp_reservations;
