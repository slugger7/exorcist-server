alter type job_type_enum add value 'scan_library'; -- triggers scan path jobs for each path in the library
alter type job_type_enum add value 'refresh_metadata'; -- triggers a refresh metadata job on one media entity
alter type job_type_enum add value 'refresh_library_metadata'; -- triggers refresh_metadata jobs on all entities that exist in the library
alter type job_type_enum add value 'generate_chapters'; -- triggers generate_thumbnail jobs for a media entity
alter type job_type_enum add value 'generate_library_chapters'; -- triggerrs generate_chapters job for all entities in the library
