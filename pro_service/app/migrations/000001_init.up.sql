CREATE TABLE stand
(
    id    serial       not null unique,
    title varchar(255) not null,
);

CREATE TABLE flange
(
    id    serial       not null unique,
    title varchar(255) not null,
    short varchar(255) not null,
);

CREATE TABLE additional
(
    id          serial not null unique,
    materials   text
    mod         text
    temperature text
    mounting    text
    graphite    text
    type_fl     text
);

CREATE TABLE snp
(
    id          serial                                      not null unique,
    stand_id    int references stand (id) on delete cascade not null,
    type_p      text
    type_fl     text
    filler      text
    materials   text
    mod         text
    temperature text
    mounting    text
    graphite    text
)