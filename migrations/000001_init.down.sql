begin;
    drop table if exists user_email_confirmations;
    drop table if exists user_password_resets;
    drop trigger remove_confirmations_after_confirmation on users;
    drop function confirm_email();
    drop table if exists users;
commit;
