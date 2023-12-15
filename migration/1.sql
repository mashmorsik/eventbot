create table if not exists public.users
(
    user_id bigint not null primary key
);

create table if not exists public.events
(
    id        serial primary key,
    user_id   bigint                   not null references users (user_id),
    name      text                     not null,
    check ( name <> '' ),
    time_date timestamp with time zone not null,
    weekly    bool,
    monthly   bool,
    yearly    bool
-- check only one is true
);

-- drop table public.users;
-- drop table public.events

