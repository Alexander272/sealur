CREATE TABLE stand (
    "id"  serial not null unique,
    "title" varchar(255) not null
);

CREATE TABLE flange (
    "id" serial not null unique,
    "title" varchar(255) not null,
    "short" varchar(255) not null
);

CREATE TABLE st_fl (
    "id" serial not null unique,
    "stand_id" int references stand (id) on delete cascade not null,
    "fl_ids" int references flange (id) on delete cascade not null
);

CREATE TABLE type_fl (
    "id" serial not null unique,
    "title" varchar(255) not null,
    "desc" varchar(255),
    "short" varchar(255) not null,
    "basis" boolean
);

CREATE TABLE additional (
    "id" serial not null unique,
    "materials" text,
    "mod" text,
    "temperature" text,
    "mounting" text,
    "graphite" text
);

CREATE TABLE snp (
    "id" serial not null unique,
    "stand_id" int references stand (id) on delete cascade not null,
    "type_pr" text,
    "type_fl_id" int references type_fl (id) on delete cascade not null,
    "filler" text,
    "materials" text,
    "mod" text,
    "temperature" text,
    "mounting" text,
    "graphite" text
);