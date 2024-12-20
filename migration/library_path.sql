create table library_path (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  library_id uuid not null,
  path varchar not null unique,
  created timestamp default current_timestamp not null,
  modified timestamp default current_timestamp not null,
  constraint fk_library
    foreign key(library_id) references library(id)
    on delete cascade
)