create table users (id serial not null , login varchar(20) primary key not null , password varchar not null , balance int not null );
create type order_status as enum ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
create table orders (id serial primary key, user_login varchar(20) references users(login), status order_status, accrual int, uploaded_at date);