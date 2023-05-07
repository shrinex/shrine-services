-- goctl model mysql ddl -src="./proto/sql/brand.sql" -dir="./proto/model" -c --style=go_zero

create database if not exists product;
use product;

create table brand
(
    brand_id    bigint       not null comment '品牌ID',
    name        varchar(24)  not null comment '品牌名称',
    remark      varchar(64)  not null comment '品牌描述',
    logo        varchar(128) not null comment '品牌图片',
    status      tinyint      not null comment '状态:0-禁用,1-启用,2-已删除',
    weight      int          not null comment '排序',
    create_time datetime     not null default current_timestamp comment '创建时间',
    update_time datetime     not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (brand_id),
    unique key uk_name (name)
) comment ='品牌表';