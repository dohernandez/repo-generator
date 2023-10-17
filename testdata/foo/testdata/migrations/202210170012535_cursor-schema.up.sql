create table if not exists "cursor" (
    id uuid not null,
    name varchar(255) not null,
    position uuid not null,
    leader uuid null,
    leader_elected_at timestamp null,
    created_at timestamp not null,
    updated_at timestamp not null,

    primary key (id)
);

create unique index unique_cursor_name on "cursor" (name);