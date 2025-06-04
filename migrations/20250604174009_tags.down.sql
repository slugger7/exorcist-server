alter table media_tag drop constraint fk_media_tag_media;
alter table media_tag drop constraint fk_media_tag_tag;
drop table media_tag;

alter table tag_alias drop constraint fk_tag_alias_tag;
drop table tag_alias;

drop table tag;
