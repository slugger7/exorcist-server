create type job_type_enum as enum ('update_existing_videos', 'scan_path','generate_checksum');
create type job_status_enum as enum ('not_started', 'in_progress', 'failed', 'completed', 'cancelled');

create table job (
    id uuid primary key default gen_random_uuid(),
    job_type job_type_enum not null,
    status job_status_enum not null,
    data jsonb,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null
)
