begin;

  alter type job_type_enum rename to old_job_type_enum;

  create type job_type_enum as enum ('update_existing_videos', 'scan_path','generate_checksum', 'generate_thumbnail');

  alter table job rename column job_type to old_job_type;

  alter table job add job_type job_type_enum not null default 'scan_path';

  update job set job_type = old_job_type::text::job_type_enum;

  alter table job drop column old_job_type;

  drop type old_job_type_enum;

commit;
