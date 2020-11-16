CREATE SCHEMA distate;
SET search_path = distate, public;

GRANT USAGE ON SCHEMA distate TO distate_user;

---- create above / drop below ----

DROP SCHEMA IF EXISTS distate CASCADE;
