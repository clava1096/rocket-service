-- +goose Up
create type order_status as enum('PENDING_PAYMENT', 'PAID', 'CANCELLED');

create type order_payment_method as enum('UNKNOWN', 'CARD', 'SBP', 'CREDIT_CARD', 'INVESTOR_MONEY');

create table orders(
    uuid    uuid primary key,
    user_uuid   uuid not null ,
    part_uuids  uuid[] not null,
    total_price double precision not null ,
    status order_status not null,
    transaction_uuid text,
    payment_method order_payment_method,
    created_at timestamp not null default now(),
    updated_at timestamp
);

-- +goose Down
drop table orders;
drop type order_status;
drop type order_payment_method;