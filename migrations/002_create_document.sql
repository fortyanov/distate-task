CREATE TABLE distate.document
(
    id     serial PRIMARY KEY,
    name   varchar(256),
    date   timestamptz,
    number bigint NOT NULL,
    sum    varchar(256)
);

---- create above / drop below ----

DROP TABLE distate.document;
