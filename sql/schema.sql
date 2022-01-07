create table if not exists customers
(
    id bigserial primary key,
    name	text not null,
    phone 	text 	not null unique,
    password text 	not null,
    active 	boolean not null default true,
    created timestamp not null default current_timestamp 
);