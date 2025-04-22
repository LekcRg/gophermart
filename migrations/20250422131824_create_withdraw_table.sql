-- +goose Up
CREATE TABLE withdraw (
  id SERIAL PRIMARY KEY,
  order_id varchar(50) NOT NULL,
  sum DOUBLE PRECISION DEFAULT 0,
  user_login varchar(30) NOT NULL REFERENCES users (login),
  processed_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE withdraw;
