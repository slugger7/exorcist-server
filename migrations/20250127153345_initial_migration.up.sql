begin;
  create type library_type_enum as enum ('image', 'video', 'mixed');

  create table library
  (
    id uuid primary key default gen_random_uuid(),
    "name" varchar not null unique,
    library_type library_type_enum default 'mixed' not null,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null
  );

  create table library_path
  (
    id uuid primary key default gen_random_uuid(),
    library_id uuid not null,
    path varchar not null unique,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null,
    constraint fk_library_path_library
      foreign key(library_id) references library(id)
      on delete cascade
  );

  create type media_type_enum as enum ('primary', 'asset');

  create table media
  (
    id uuid primary key default gen_random_uuid(),
    library_path_id uuid not null,
    path varchar not null,
    title varchar not null,
    media_type media_type_enum default 'primary' not null,
    size bigint not null,
    checksum char(32),
    added timestamp default current_timestamp not null,
    deleted boolean default false not null,
    exists boolean default true not null,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null,
    constraint fk_media_library_path 
      foreign key (library_path_id) 
      references library_path (id) 
      on delete cascade
  );

  create table video
  (
    id uuid primary key default gen_random_uuid(),
    media_id uuid not null,
    height int not null,
    width int not null,
    runtime double precision not null,
    constraint fk_video_media 
      foreign key (media_id) 
      references media(id) 
      on delete cascade
  );

  create table "image"
  (
    id uuid primary key default gen_random_uuid(),
    media_id uuid not null,
    height int not null,
    width int not null,
    constraint fk_image_media 
      foreign key (media_id) 
      references media(id) 
      on delete cascade
  );

  create type media_relation_type_enum as enum
  ('thumbnail', 'chapter', 'media');

  create table media_relation
  (
    id uuid primary key default gen_random_uuid(),
    media_id uuid not null,
    related_to uuid not null, 
    relation_type media_relation_type_enum not null,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null,
    constraint fk_media_relation_media 
      foreign key(media_id) 
      references "media"(id) 
      on delete cascade,
    constraint fk_media_relation_related_to 
      foreign key(related_to) 
      references "media"(id) 
      on delete cascade
  );

  create table "user"
  (
    id uuid primary key default gen_random_uuid(),
    username varchar not null,
    password varchar not null,
    active boolean default true not null,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null
  );

  create unique index unique_username_constraint on "user"(username);

  insert into "user"
    (username, password)
  values
    ('admin', '$2a$10$5nmh/cOu.dzk05V7lfBqQua9FO6nG.aQTGTJQFB26DGMSMwp5FWxu');
  -- password = admin

  create type job_type_enum as enum
  ('update_existing_videos', 'scan_path','generate_checksum', 'generate_thumbnail');
  create type job_status_enum as enum
  ('not_started', 'in_progress', 'failed', 'completed', 'cancelled');

  create table job
  (
    id uuid primary key default gen_random_uuid(),
    parent uuid,
    priority smallint default 2 not null,
    job_type job_type_enum not null,
    status job_status_enum not null,
    data jsonb,
    outcome varchar,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null
  );
commit;
