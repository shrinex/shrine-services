-- goctl model mysql ddl -src="./proto/sql/category.sql" -dir="./proto/model" -c --style=go_zero

create database if not exists product;
use product;

drop table if exists category;
create table category
(
    category_id bigint       not null comment '分类ID',
    parent_id   bigint       not null comment '父ID',
    group_id    bigint       not null comment '分类分组ID(仅1级分类有效)',
    name        varchar(24)  not null comment '分类名称',
    remark      varchar(36)  not null comment '分类描述',
    icon        varchar(128) not null comment '分类图标',
    level       tinyint      not null comment '层级,从1开始',
    status      tinyint      not null comment '状态:0-禁用,1-启用,2-已删除',
    weight      int          not null comment '排序',
    create_time datetime     not null default current_timestamp comment '创建时间',
    update_time datetime     not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (category_id),
    unique key uk_name (name),
    key idx_parentid (parent_id)
) comment '产品分类表';