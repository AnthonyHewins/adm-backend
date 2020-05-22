begin;

drop table user_password_resets;

create table if not exists user_password_resets(
       user_id        bigint references users(id) on delete cascade unique not null,
       token          char(40) unique not null,
       reset_at       timestamp with time zone not null default CURRENT_TIMESTAMP
);

commit;
