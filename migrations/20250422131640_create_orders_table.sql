-- +goose Up
CREATE TYPE order_status AS ENUM (
  'NEW',
  'PROCESSING',
  'INVALID',
  'PROCESSED'
);

CREATE TABLE orders (
  id serial NOT NULL PRIMARY KEY,
  order_id varchar(50) NOT NULL UNIQUE,
  status order_status NOT NULL,
  accrual double precision,
  uploaded_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  user_login varchar(30) NOT NULL REFERENCES users (LOGIN)
);

-- +goose Down
DROP TABLE orders;
DROP TYPE order_status;
