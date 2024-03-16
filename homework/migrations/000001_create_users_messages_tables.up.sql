CREATE TABLE IF NOT EXISTS users
(
    nickname      varchar(64) primary key,
    password_hash varchar(256) not null
);

CREATE TABLE IF NOT EXISTS messages
(
    id        serial primary key,
    sent_time timestamptz not null,
    nickname  varchar(64) not null references users (nickname),
    message   text        not null
);
