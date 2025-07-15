begin;
  create table favourite_person
  (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null,
    person_id uuid not null,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null,
    constraint fk_favourite_person_user
      foreign key(user_id)
      references "user"(id)
      on delete cascade,
    constraint fk_favourite_person_person
      foreign key (person_id)
      references "person"(id)
      on delete cascade
  );

  create table favourite_media
  (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null,
    media_id uuid not null,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null,
    constraint fk_favourite_media_user
      foreign key(user_id)
      references "user"(id)
      on delete cascade,
    constraint fk_favourite_media_media
      foreign key(media_id)
      references "media"(id)
      on delete cascade
  );

  create table playlist
  (
    id uuid primary key default gen_random_uuid(),
    "name" varchar not null,
    user_id uuid not null,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null,
    constraint fk_playlist_user
      foreign key(user_id)
      references "user"(id)
      on delete cascade
  );

  create table playlist_media
  (
    id uuid primary key default gen_random_uuid(),
    playlist_id uuid not null,
    media_id uuid not null,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null,
    constraint fk_playlist_media_playlist
      foreign key(playlist_id)
      references "playlist"(id)
      on delete cascade,
    constraint fk_playlist_media_media
      foreign key(media_id)
      references "media"(id)
  );

  create table media_progress
  (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null,
    media_id uuid not null,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null,
    constraint fk_media_progress_user
      foreign key(user_id)
      references "user"(id)
      on delete cascade,
    constraint fk_media_progress_media
      foreign key(media_id)
      references "media"(id)
      on delete cascade
  );
commit;
