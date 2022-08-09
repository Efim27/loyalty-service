CREATE TABLE IF NOT EXISTS "user"
(
    id       serial       NOT NULL,
    login    varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    balance  float        NOT NULL,
    CONSTRAINT user_pk PRIMARY KEY (id)
);