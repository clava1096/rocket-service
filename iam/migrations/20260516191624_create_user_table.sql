-- +goose Up
create table users (
    uuid uuid primary key default uuid_generate_v4(),
    login  varchar(100) not null unique (login),
    email varchar(255) not null unique (email),
    password_hash varchar(255) not null,
    notification_methods notification_method,
    createdAt timestamp not null default now(),
    updatedAt timestamp
);

create table notification_method (
    id serial primary key,
    user_uuid uuid references users(uuid) on delete cascade
    provider_name varchar(50) not null,
    target text not null
);

-- +goose Down
drop table users;
drop table notification_method;