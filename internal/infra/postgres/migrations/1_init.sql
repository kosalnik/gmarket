-- +goose Up
CREATE TYPE order_status AS ENUM (
    'NEW',
    'PROCESSING',
    'INVALID',
    'PROCESSED'
);

CREATE TYPE withdraw_status AS ENUM (
    'REGISTERED',
    'INVALID',
    'PROCESSING',
    'PROCESSED'
);

CREATE TABLE "user" (
    login character varying(255) NOT NULL PRIMARY KEY,
    password char(32) NOT NULL,
    balance numeric(20, 2) NOT NULL DEFAULT 0
);

CREATE TABLE "order" (
    number numeric(200) NOT NULL PRIMARY KEY,
    status order_status NOT NULL DEFAULT 'NEW',
    accrual numeric(20, 2) NOT NULL DEFAULT 0,
    uploaded_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE "withdraw" (
    "order" numeric(200) NOT NULL PRIMARY KEY,
    sum numeric(20, 2) NOT NULL DEFAULT 0,
    status withdraw_status NOT NULL DEFAULT 'REGISTERED'
);

-- +goose Down
DROP TABLE "withdraw";
DROP TABLE "order";
DROP TABLE "user";