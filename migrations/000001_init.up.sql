begin;

create table if not exists users (
       id               bigserial primary key not null,
       email            varchar(255) not null unique,
       password         varchar(60) not null,
       registered_at    timestamp with time zone not null default CURRENT_TIMESTAMP,
       confirmed_at     timestamp,

       constraint register_before_confirmation check (registered_at <= confirmed_at)
);

create table if not exists user_email_confirmations (
       user_id  bigint references users(id) on delete cascade unique not null,
       token    char(40) unique not null
);

create table if not exists user_password_resets(
       user_id        bigint references users(id) on delete cascade unique not null,
       token          char(40) unique not null,
       reset_at       timestamp with time zone not null default CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION confirm_email() RETURNS TRIGGER AS
'
BEGIN
delete from user_email_confirmations where user_id = new.id;
RETURN NEW;
END;
' LANGUAGE plpgsql;

create trigger remove_confirmations_after_confirmation
       after update of confirmed_at on users
       for each row
       when (new.confirmed_at is not null)
       execute function confirm_email();

commit;
