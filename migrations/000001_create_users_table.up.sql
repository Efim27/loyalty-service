CREATE TABLE IF NOT EXISTS "user"
(
    id       serial       NOT NULL,
    login    varchar(255) NOT NULL,
    password varchar(255) NOT NULL UNIQUE,
    balance  float        NOT NULL,
    CONSTRAINT user_pk PRIMARY KEY (id)
);