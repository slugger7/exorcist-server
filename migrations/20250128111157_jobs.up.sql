create type job_type as enum ('update_existing_videos', 'scan_path','generate_checksum');
create type job_status as enum ('not_started', 'in_progress', 'failed', 'completed', 'cancelled');

create table job (
    id uuid primary key default gen_random_uuid(),
    job_type job_type not null,
    status job_status not null,
    data jsonb
)
