create table library (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "name" varchar not null unique,
  created timestamp default current_timestamp not null,
  modified timestamp default current_timestamp not null
)