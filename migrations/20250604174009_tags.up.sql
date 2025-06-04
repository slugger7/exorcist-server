create table tag
(
  id uuid primary key default gen_random_uuid(),
  name varchar not null unique,
  created timestamp default current_timestamp not null,
  modified timestamp default current_timestamp not null
);

create table tag_alias
(
  id uuid primary key default gen_random_uuid(),
  alias varchar not null,
  created timestamp default current_timestamp not null,
  modified timestamp default current_timestamp not null,
  tag_id uuid not null,
  constraint fk_tag_alias_tag
    foreign key(tag_id) references tag(id)
    on delete cascade
);

create table media_tag
(
  id uuid primary key default gen_random_uuid(),
  created timestamp default current_timestamp not null,
  modified timestamp default current_timestamp not null,
  media_id uuid not null,
  tag_id uuid not null,
  constraint fk_media_tag_media
    foreign key(media_id) references media(id)
    on delete cascade,
  constraint fk_media_tag_tag
    foreign key(tag_id) references tag(id)
    on delete cascade
);

