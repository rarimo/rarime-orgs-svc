-- +migrate Up
create table claims_schemas
(
    id uuid primary key default gen_random_uuid(),
    action_type text not null,
    schema_type text not null,
    schema_url text not null,
    created_at timestamp without time zone not null default now(),
    updated_at timestamp without time zone not null default now()
);

-- +migrate Down
drop table claims_schemas;
