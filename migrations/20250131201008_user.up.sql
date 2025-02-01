create table users (
  id uuid primary key default gen_random_uuid(),
  username varchar not null, 
  password varchar not null,
  active boolean default true not null,
  created timestamp default current_timestamp not null,
  modified timestamp default current_timestamp not null
);

create unique index unique_username_constraint on users(username);
