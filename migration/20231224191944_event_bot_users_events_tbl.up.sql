create table if not exists public.users
(
    user_id bigint not null primary key
);

create table if not exists public.events
(
    id        serial primary key,
    user_id   bigint                   not null references users (user_id),
    chat_id   bigint                   not null,
    name      text                     not null,
    check ( name <> '' ),
    time_date timestamp with time zone not null,
    cron      text                     not null,
    last_fired     timestamp with time zone,
    disabled    bool                    default false
);

create table if not exists public.cron
(
    id      int primary key,
    cron    json                       not null
)
