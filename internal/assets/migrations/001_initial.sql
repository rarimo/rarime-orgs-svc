-- +migrate Up
drop function if exists trigger_set_updated_at cascade;
CREATE FUNCTION trigger_set_updated_at() RETURNS trigger
    LANGUAGE plpgsql
AS $$ BEGIN NEW.updated_at = NOW() at time zone 'utc'; RETURN NEW; END; $$;

create table if not exists users
(
    id         uuid primary key                     default gen_random_uuid(),
    did        text                        not null,
    org_id     uuid                        not null,
    role       smallint                    not null default 0,
    created_at timestamp without time zone not null default now(),
    updated_at timestamp without time zone not null default now()
);

create index if not exists users_did_index on users using btree (did);

create trigger set_updated_at
    before update
    on users
    for each row
execute function trigger_set_updated_at();

create table if not exists organizations
(
    id                  uuid primary key                     default gen_random_uuid(),
    did                 text,
    owner               uuid                        not null,
    domain              text                        not null,
    metadata            jsonb                       not null default '{}'::jsonb,
    status              smallint                    not null default 0,
    verification_code   text,
    issued_claims_count bigint                      not null default 0,
    members_count       int                         not null default 0,
    created_at          timestamp without time zone not null default now(),
    updated_at          timestamp without time zone not null default now()
);

create index if not exists organizations_owner_index on organizations using btree (owner);

create trigger set_updated_at
    before update
    on organizations
    for each row
execute function trigger_set_updated_at();

create table if not exists groups
(
    id         uuid primary key                     default gen_random_uuid(),
    org_id     uuid                        not null,
    metadata   jsonb                       not null default '{}'::jsonb,
    rules      jsonb                       not null default '{}'::jsonb,
    created_at timestamp without time zone not null default now()
);

create index if not exists groups_organization_index on groups using btree (org_id);

create table if not exists group_users
(
    id         uuid primary key                     default gen_random_uuid(),
    group_id   uuid                        not null,
    user_id    uuid                        not null,
    role       smallint                    not null default 0,
    created_at timestamp without time zone not null default now(),
    updated_at timestamp without time zone not null default now()
);

create index if not exists group_users_group_index on group_users using btree (group_id);
create index if not exists group_users_user_index on group_users using btree (user_id);

create trigger set_updated_at
    before update
    on group_users
    for each row
execute function trigger_set_updated_at();

create table if not exists requests
(
    id         uuid primary key                     default gen_random_uuid(),
    org_id     uuid                        not null,
    group_id   uuid                        not null,
    user_id    uuid,
    metadata   jsonb                       not null default '{}'::jsonb,
    status     smallint                    not null default 0,
    created_at timestamp without time zone not null default now(),
    updated_at timestamp without time zone not null default now()
);

create index if not exists requests_organization_index on requests using btree (org_id);
create index if not exists requests_group_index on requests using btree (group_id);
create index if not exists requests_user_index on requests using btree (user_id);

create trigger set_updated_at
    before update
    on requests
    for each row
execute function trigger_set_updated_at();

create table if not exists email_invitations
(
    id         uuid primary key                     default gen_random_uuid(),
    req_id     uuid                        not null,
    org_id     uuid                        not null,
    group_id   uuid                        not null,
    email      text                        not null,
    otp        text                        not null,
    created_at timestamp without time zone not null default now()
);

create index if not exists email_invitations_organization_index on email_invitations using btree (org_id);
create index if not exists email_invitations_group_index on email_invitations using btree (group_id);
create index if not exists email_invitations_request_index on email_invitations using btree (req_id);

-- +migrate Down
drop index if exists email_invitations_request_index;
drop index if exists email_invitations_group_index;
drop index if exists email_invitations_organization_index;
drop table if exists email_invitations;

drop trigger if exists set_updated_at on requests;
drop index if exists requests_user_index;
drop index if exists requests_group_index;
drop index if exists requests_organization_index;
drop table if exists requests;

drop trigger if exists set_updated_at on group_users;
drop index if exists group_users_user_index;
drop index if exists group_users_group_index;
drop table if exists group_users;

drop index if exists groups_organization_index;
drop table if exists groups;

drop trigger if exists set_updated_at on organizations;
drop index if exists organizations_owner_index;
drop table if exists organizations;

drop trigger if exists set_updated_at on users;
drop index if exists users_did_index;
drop table if exists users;

drop function if exists trigger_set_updated_at;
