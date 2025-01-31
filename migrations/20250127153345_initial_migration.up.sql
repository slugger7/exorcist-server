create table library (
  id uuid primary key default gen_random_uuid(),
  "name" varchar not null unique,
  created timestamp default current_timestamp not null,
  modified timestamp default current_timestamp not null
);

create table library_path (
  id uuid primary key default gen_random_uuid(),
  library_id uuid not null,
  path varchar not null unique,
  created timestamp default current_timestamp not null,
  modified timestamp default current_timestamp not null,
  constraint fk_library
    foreign key(library_id) references library(id)
    on delete cascade
);

create table video (
  id uuid primary key default gen_random_uuid(),
  library_path_id uuid not null,
  relative_path varchar not null,
  title varchar not null,
  file_name varchar not null,
  height int not null,
  width int not null,
  runtime bigint not null,
  size bigint not null,
  checksum char(40),
  added timestamp default current_timestamp not null,
  deleted boolean default false not null,
  exists boolean default true not null,
  created timestamp default current_timestamp not null,
  modified timestamp default current_timestamp not null,
  constraint fk_library_path
    foreign key(library_path_id) references library_path(id)
    on delete cascade
);
