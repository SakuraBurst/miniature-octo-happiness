create table if not exists users (id serial not null , login varchar(20) primary key not null unique , password varchar not null , balance real not null default 0.0, withdraw real not null default 0.0);
DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
            create type order_status as enum ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
        END IF;
    END
$$;
create table if not exists orders (id serial not null , user_login varchar(20) references users(login) not null , order_id varchar primary key not null ,  status order_status not null , accrual real not null default 0.0, uploaded_at timestamp not null );
create table if not exists withdraws (id serial not null, user_login varchar(20) references users(login) not null, order_id varchar not null ,sum real not null , processed_at timestamp not null);