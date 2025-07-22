alter table library drop constraint library_unique_ghost_id;
alter table library drop column ghost_id;

alter table library_path drop constraint library_path_unique_ghost_id;
alter table library_path drop column ghost_id;

alter table media drop constraint media_unique_ghost_id;
alter table media drop column ghost_id;

alter table video drop constraint video_unique_ghost_id;
alter table video drop column ghost_id;

alter table media_relation drop constraint media_relation_unique_ghost_id;
alter table media_relation drop column ghost_id;

alter table "user" drop constraint user_unique_ghost_id;
alter table "user" drop column ghost_id;

alter table person drop constraint person_unique_ghost_id;
alter table person drop column ghost_id;

alter table media_person drop constraint media_person_unique_ghost_id;
alter table media_person drop column ghost_id;

alter table tag drop constraint tag_unique_ghost_id;
alter table tag drop column ghost_id;
 
alter table media_tag drop constraint media_tag_unique_ghost_id;
alter table media_tag drop column ghost_id;

alter table favourite_person drop constraint favourite_person_unique_ghost_id;
alter table favourite_person drop column ghost_id;

alter table favourite_media drop constraint favourite_media_unique_ghost_id;
alter table favourite_media drop column ghost_id;

alter table playlist drop constraint playlist_unique_ghost_id;
alter table playlist drop column ghost_id;

alter table playlist_media drop constraint playlist_media_unique_ghost_id;
alter table playlist_media drop column ghost_id;

alter table media_progress drop constraint media_progress_unique_ghost_id;
alter table media_progress drop column ghost_id;
