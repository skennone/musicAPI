CREATE TABLE IF NOT EXISTS songs (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone not null default now(),
    title text not null,
    artist text not null,
    year integer not null,
    length integer not null,
    genres text[] not null,
    version integer not null default 1
);
