alter table video drop constraint fk_video_library_path;
drop table if exists video;

alter table library_path drop constraint fk_library_path_library;
drop table if exists library_path;

drop table if exists library;
