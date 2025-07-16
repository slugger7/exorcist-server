alter table library add column ghost_id integer null;
alter table library add constraint library_unique_ghost_id unique (ghost_id);

alter table library_path add column ghost_id integer;
alter table library_path add constraint library_path_unique_ghost_id unique (ghost_id);

alter table video add column ghost_id integer;
alter table video add constraint video_unique_ghost_id unique (ghost_id);

alter table media_relation add column ghost_id integer;
alter table media_relation add constraint media_relation_unique_ghost_id unique (ghost_id);

alter table "user" add column ghost_id integer;
alter table "user" add constraint user_unique_ghost_id unique (ghost_id);

alter table person add column ghost_id integer;
alter table person add constraint person_unique_ghost_id unique (ghost_id);

alter table media_person add column ghost_id integer;
alter table media_person add constraint media_person_unique_ghost_id unique (ghost_id);

alter table tag add column ghost_id integer;
alter table tag add constraint tag_unique_ghost_id unique (ghost_id);

alter table media_tag add column ghost_id integer;
alter table media_tag add constraint media_tag_unique_ghost_id unique (ghost_id);

alter table favourite_person add column ghost_id integer;
alter table favourite_person add constraint favourite_person_unique_ghost_id unique (ghost_id);

alter table favourite_media add column ghost_id integer;
alter table favourite_media add constraint favourite_media_unique_ghost_id unique (ghost_id);

alter table playlist add column ghost_id integer;
alter table playlist add constraint playlist_unique_ghost_id unique (ghost_id);

alter table playlist_media add column ghost_id integer;
alter table playlist_media add constraint playlist_media_unique_ghost_id unique (ghost_id);

alter table media_progress add column ghost_id integer;
alter table media_progress add constraint media_progress_unique_ghost_id unique (ghost_id);
