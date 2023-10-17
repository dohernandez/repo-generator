create table if not exists "block" (
    id uuid not null,
    chain_id int not null,
    number bigint not null,
    hash varchar(66),
    parent_hash varchar(66),
    block_timestamp timestamp,

    primary key(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_block_hash" ON "block" (hash);