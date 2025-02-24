create table "image" (
    id uuid primary key default gen_random_uuid(),
    name varchar not null,
    path varchar not null,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null
);

create unique index unique_image_path on "image"(path);

create type video_image_type_enum as enum ('thumbnail', 'chapter');

create table "video_image" (
    id uuid primary key default gen_random_uuid(),
    video_id uuid not null,
    image_id uuid not null,
    video_image_type video_image_type_enum not null,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null,
    constraint fk_video_image_video foreign key(video_id) references "video"(id) on delete cascade,
    constraint fk_video_image_image foreign key(image_id) references "image"(id) on delete cascade
);
