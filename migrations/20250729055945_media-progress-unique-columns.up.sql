alter table media_progress add constraint uq_media_progress_media_id_user_id unique (media_id, user_id)
