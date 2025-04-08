CREATE TABLE IF NOT EXISTS users (
                                     id serial primary key,
                                     telegram_id bigint not null unique,
                                     channel_id bigint not null,
                                     created_at timestamp default now()
    );