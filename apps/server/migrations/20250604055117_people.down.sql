alter table media_person drop constraint fk_media_person_media;
alter table media_person drop constraint fk_media_person_person;
drop table media_person;

alter table person_alias drop constraint fk_person_alias_person;
drop table person_alias;

drop table person;
