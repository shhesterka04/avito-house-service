-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_type AS ENUM ('client', 'moderator');
CREATE TYPE flat_status AS ENUM ('created', 'on moderation', 'approved', 'declined');

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    type user_type NOT NULL
);

CREATE TABLE house
(
    id SERIAL PRIMARY KEY,
    address VARCHAR(255) NOT NULL,
    year INT NOT NULL,
    developer VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE flats
(
    id SERIAL PRIMARY KEY,
    house_id INT NOT NULL REFERENCES house(id),
    status flat_status NOT NULL,
    number INT NOT NULL,
    rooms INT NOT NULL,
    price INT NOT NULL
);

CREATE INDEX idx_house_id ON flats(house_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE flats;
DROP TABLE house;

DROP TYPE user_type;
DROP TYPE flat_status;

DROP INDEX idx_house_id;
-- +goose StatementEnd
