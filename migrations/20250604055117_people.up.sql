create table person
(
  id uuid primary key default gen_random_uuid(),
  name varchar not null unique,
  created timestamp default current_timestamp not null,
  modified timestamp default current_timestamp not null
);

create table person_alias
(
  id uuid primary key default gen_random_uuid(),
  alias varchar not null,
  created timestamp default current_timestamp not null,
  modified timestamp default current_timestamp not null,
  person_id uuid not null,
  constraint fk_person_alias_person
    foreign key(person_id) references person(id)
    on delete cascade
);

create table media_person
(
  id uuid primary key default gen_random_uuid(),
  created timestamp default current_timestamp not null,
  modified timestamp default current_timestamp not null,
  media_id uuid not null,
  person_id uuid not null,
  constraint fk_media_person_media
    foreign key(media_id) references media(id)
    on delete cascade,
  constraint fk_media_person_person
    foreign key(person_id) references person(id)
    on delete cascade
);
