-- goctl model mysql ddl -src="./unit/product/proto/sql/category_group.sql" -dir="./unit/product/proto/model" -c --style=go_zero

create database if not exists product;
use product;

drop table if exists category_group;
create table category_group
(
    group_id    bigint       not null comment '分类分组ID',
    name        varchar(24)  not null comment '分类分组名称',
    icon        varchar(128) not null comment '分类分组图标',
    status      tinyint      not null comment '状态:0-禁用,1-启用',
    weight      int          not null comment '排序',
    create_time datetime     not null default current_timestamp comment '创建时间',
    update_time datetime     not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (group_id),
    unique key uk_name (name)
) comment '一级分类分组表';