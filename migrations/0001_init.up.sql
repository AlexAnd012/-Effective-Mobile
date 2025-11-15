create extension if not exists pgcrypto;

create table if not exists subscriptions (
id uuid primary key default gen_random_uuid(),
service_name text not null,
price int not null check (price > 0),
user_id uuid not null,
start_date date not null,
end_date date null,
check (end_date is null or end_date >= start_date)
);

create index if not exists ix_subs_user on subscriptions(user_id);
create index if not exists ix_subs_service on subscriptions(service_name);
