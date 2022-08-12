CREATE TABLE IF NOT EXISTS "order"
(
    id         serial                   NOT NULL,
    number     bigint                   NOT NULL UNIQUE,
    status     varchar(255)             NOT NULL CHECK (status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED')),
    accrual    FLOAT                    NOT NULL CHECK (accrual >= 0),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT order_pk PRIMARY KEY (id)
);