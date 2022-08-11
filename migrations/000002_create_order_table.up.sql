CREATE TABLE IF NOT EXISTS "order"
(
    id         serial       NOT NULL,
    number     bigint       NOT NULL,
    status     varchar(255) NOT NULL CHECK (status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED')),
    accrual    FLOAT        NOT NULL CHECK (accrual >= 0),
    created_at TIMESTAMP    NOT NULL,
    CONSTRAINT order_pk PRIMARY KEY (id)
);