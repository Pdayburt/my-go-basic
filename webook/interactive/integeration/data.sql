create database if not exists webook;
create table if not exists webook.interactives
(
    id          bigint auto_increment
        primary key,
    biz_id      bigint       null,
    biz         varchar(128) null,
    read_cnt    bigint       null,
    collect_cnt bigint       null,
    like_cnt    bigint       null,
    ctime       bigint       null,
    utime       bigint       null,
    constraint biz_type_id
        unique (biz_id, biz)
);

create table if not exists webook.user_collection_bizs
(
    id     bigint auto_increment
        primary key,
    cid    bigint       null,
    biz_id bigint       null,
    biz    varchar(128) null,
    uid    bigint       null,
    ctime  bigint       null,
    utime  bigint       null,
    constraint biz_type_id_uid
        unique (biz_id, biz, uid)
);

create index idx_user_collection_bizs_cid
    on webook.user_collection_bizs (cid);

create table if not exists webook.user_like_bizs
(
    id     bigint auto_increment
        primary key,
    biz_id bigint           null,
    biz    varchar(128)     null,
    uid    bigint           null,
    status tinyint unsigned null,
    ctime  bigint           null,
    utime  bigint           null,
    constraint biz_type_id_uid
        unique (biz_id, biz, uid)
);

INSERT INTO `interactives`(`biz_id`, `biz`, `read_cnt`, `collect_cnt`, `like_cnt`, `ctime`, `utime`)
VALUES(1,"test",1494,3656,7103,1728383366915,1728383366915),
(2,"test",4935,7185,1621,1728383366915,1728383366915),
(3,"test",4899,5955,8722,1728383366915,1728383366915),
(4,"test",8720,9982,4892,1728383366915,1728383366915),
(5,"test",6889,6173,6594,1728383366915,1728383366915),
(6,"test",2388,2158,5186,1728383366915,1728383366915),
(7,"test",6999,2697,6345,1728383366915,1728383366915),
(8,"test",9579,1292,4372,1728383366915,1728383366915),
(9,"test",4234,9594,258,1728383366915,1728383366915),
(10,"test",8810,1587,351,1728383366915,1728383366915)