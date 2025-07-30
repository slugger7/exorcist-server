with duplicate_progress as (select id,
        row_number() over (partition by user_id, media_id order by timestamp desc) as rn
from media_progress)
delete from media_progress
where id in (
  select id
  from duplicate_progress 
  where duplicate_progress.rn > 1
);

alter table media_progress add constraint uq_media_progress_media_id_user_id unique (media_id, user_id);
