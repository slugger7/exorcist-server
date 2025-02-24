alter table "video_image" drop constraint fk_video_image_video;
alter table "video_image" drop constraint fk_video_image_image;

drop table if exists "video_image";
drop table if exists "image";

drop type if exists video_image_type_enum;
