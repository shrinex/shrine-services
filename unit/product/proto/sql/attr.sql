-- goctl model mysql ddl -src="./unit/product/proto/sql/attr.sql" -dir="./unit/product/proto/model" -c --style=go_zero

create database if not exists product;
use product;

drop table if exists attr;
create table attr
(
    attr_id      bigint      not null comment '属性ID',
    name         varchar(20) not null comment '属性名称',
    remark       varchar(36) not null comment '属性描述',
    type         tinyint     not null comment '类型:0-规格,1-属性',
    customizable tinyint     not null comment '属性值是否支持自定义:0-否,1-是',
    create_time  datetime    not null default current_timestamp comment '创建时间',
    update_time  datetime    not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (attr_id),
    unique key uk_name(name)
) comment ='属性表';