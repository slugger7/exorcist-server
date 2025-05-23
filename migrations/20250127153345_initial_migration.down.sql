alter table video_image drop constraint fk_video_image_video;
alter table video_image drop constraint fk_video_image_image;

drop table if exists video_image;

alter table image drop constraint fk_image_media;
drop table if exists image;

drop type if exists video_image_type_enum;

alter table video drop constraint fk_video_media;
drop table if exists video;

alter table media drop constraint fk_media_library_path;
drop table if exists meida;

alter table library_path drop constraint fk_library_path_library;
drop table if exists library_path;

drop table if exists library;

drop table if exists "user";

drop table if exists job;

drop type if exists job_type_enum;
drop type if exists job_status_enum;

