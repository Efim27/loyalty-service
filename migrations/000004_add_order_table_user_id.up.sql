ALTER TABLE IF EXISTS "order"
    ADD COLUMN user_id integer NOT NULL REFERENCES "user" (id);