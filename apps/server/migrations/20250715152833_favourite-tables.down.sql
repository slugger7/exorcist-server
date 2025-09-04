alter table media_progress drop constraint fk_media_progress_user;
alter table media_progress drop constraint fk_media_progress_media;
drop table if exists media_progress;

alter table playlist_media drop constraint fk_playlist_media_playlist;
alter table playlist_media drop constraint fk_playlist_media_media;
drop table if exists playlist_media;

alter table playlist drop constraint fk_playlist_user;
drop table if exists playlist;

alter table favourite_media drop constraint fk_favourite_media_user;
alter table favourite_media drop constraint fk_favourite_media_media;
drop table if exists favourite_media;

alter table favourite_person drop constraint fk_favourite_person_user;
alter table favourite_person drop constraint fk_favourite_person_person;
drop table if exists favourite_person;
