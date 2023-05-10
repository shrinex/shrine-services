-- goctl model mysql ddl -src="./unit/product/proto/sql/attr_value.sql" -dir="./unit/product/proto/model" -c --style=go_zero

create database if not exists product;
use product;

drop table if exists attr_value;
create table attr_value
(
    attr_value_id bigint      not null comment '属性值ID',
    attr_id       bigint      not null comment '属性ID',
    value         varchar(24) not null comment '属性值',
    create_time   datetime    not null default current_timestamp comment '创建时间',
    update_time   datetime    not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (attr_value_id),
    key idx_attrid (attr_id)
) comment ='属性-属性值';

