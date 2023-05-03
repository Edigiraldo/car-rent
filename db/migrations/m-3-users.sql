DROP TABLE IF EXISTS users;
CREATE TYPE USER_TYPES AS ENUM('Customer', 'Admin');
CREATE TYPE USER_STATUSES AS ENUM('Active', 'Inactive');
CREATE TABLE users (
    id uuid NOT NULL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(235) NOT NULL,
    type USER_TYPES NOT NULL,
    status USER_STATUSES NOT NULL,
    CONSTRAINT unique_email UNIQUE (email)
);
