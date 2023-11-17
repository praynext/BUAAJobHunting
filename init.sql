DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA IF NOT EXISTS public;
COMMENT ON SCHEMA public IS 'standard public schema';
SET search_path = "public";
SET TIME ZONE 'PRC';

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

create table if not exists "boss_data"
(
    id               serial primary key,
    job_name         varchar(255) not null,
    job_area         varchar(255) not null,
    salary           varchar(255) not null,
    tag_list         varchar(255) not null,
    hr_info          varchar(255) not null,
    company_logo     varchar(255) not null,
    company_name     varchar(255) not null,
    company_tag_list varchar(255) not null,
    company_url      varchar(255) not null,
    job_need         varchar(255) not null,
    job_desc         varchar(255) not null,
    job_url          varchar(255) not null,
    created_at       timestamp    not null
);

alter table "boss_data"
    owner to postgres;
