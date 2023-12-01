DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA IF NOT EXISTS public;
COMMENT ON SCHEMA public IS 'standard public schema';
SET search_path = "public";
SET TIME ZONE 'PRC';
create extension pg_trgm;

create table if not exists "user"
(
    id         serial primary key,
    name       varchar(255) not null unique,
    email      varchar(255) not null,
    phone      varchar(255),
    password   varchar(255) not null,
    created_at timestamp    not null
);

alter table "user"
    owner to postgres;

create table if not exists boss_data
(
    id               serial primary key,
    job_name         varchar(255)                 not null,
    job_area         varchar(255)                 not null,
    salary           varchar(255)                 not null,
    tag_list         varchar(255)                 not null,
    hr_info          varchar(255)                 not null,
    company_logo     varchar(255)                 not null,
    company_name     varchar(255)                 not null,
    company_tag_list varchar(255)                 not null,
    company_url      varchar(255)                 not null,
    job_need         varchar(255)                 not null,
    job_desc         varchar(255)                 not null,
    job_url          varchar(255)                 not null,
    created_at       timestamp                    not null,
    is_full          bool default random() < 0.25 not null,
    tokens           tsvector                     not null
);
create index boss_company_idx on boss_data using GIN (company_name gin_trgm_ops);
create index boss_tokens_idx on boss_data using GIN (tokens);

alter table boss_data
    owner to postgres;

create table if not exists "58_data"
(
    id           serial primary key,
    job_name     varchar(255)                 not null,
    job_area     varchar(255)                 not null,
    salary       varchar(255)                 not null,
    job_wel      varchar(255)                 not null,
    company_name varchar(255)                 not null,
    job_need     varchar(255)                 not null,
    job_url      varchar(2048)                not null,
    created_at   timestamp                    not null,
    is_full      bool default random() < 0.25 not null,
    tokens       tsvector                     not null
);
create index "58_company_idx" on "58_data" using GIN (company_name gin_trgm_ops);
create index "58_tokens_idx" on "58_data" using GIN (tokens);

alter table "58_data"
    owner to postgres;

create table message
(
    id       serial primary key,
    "from"   integer              not null,
    "to"     integer              not null,
    message  text,
    time     timestamp            not null,
    has_sent boolean default true not null
);

alter table message
    owner to postgres;

create table user_favorite_58_data
(
    id         serial primary key,
    user_id    integer not null,
    data_id    integer not null,
    created_at timestamp default now(),
    foreign key (user_id) references "user" (id)
);

alter table user_favorite_58_data
    owner to postgres;

create table user_favorite_boss_data
(
    id         serial primary key,
    user_id    integer not null,
    data_id    integer not null,
    created_at timestamp default now(),
    foreign key (user_id) references "user" (id)
);

alter table user_favorite_boss_data
    owner to postgres;

create table reminder
(
    id         serial primary key,
    user_id    integer                 not null,
    message    text,
    time       timestamp               not null,
    created_at timestamp default now(),
    has_sent   boolean   default false not null,
    foreign key (user_id) references "user" (id)
);

alter table reminder
    owner to postgres;
