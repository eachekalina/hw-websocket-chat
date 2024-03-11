CREATE TABLE messages (
    sent_time timestamptz not null,
    nickname varchar(64) not null,
    message text not null
);