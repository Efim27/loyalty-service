CREATE TABLE IF NOT EXISTS "withdrawal"
(
    id           serial    NOT NULL,
    order_number bigint    NOT NULL,
    user_id      integer   NOT NULL REFERENCES "user" (id),
    sum          FLOAT     NOT NULL CHECK (sum >= 0),
    created_at   TIMESTAMP NOT NULL,
    CONSTRAINT withdrawal_pk PRIMARY KEY (id)
);