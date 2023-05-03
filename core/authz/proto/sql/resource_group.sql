-- goctl model mysql ddl -src="./proto/sql/resource_group.sql" -dir="./proto/model" --style=go_zero

create database if not exists authz;
use authz;

drop table if exists resource_group;
create table resource_group
(
    group_id    bigint      not null comment '资源分组ID',
    name        varchar(12) not null comment '资源分组名称',
    remark      varchar(36) not null comment '备注',
    sys_type    tinyint     not null comment '系统类型: 1-平台端,2-商家端,3-普通用户',
    create_time datetime    not null default current_timestamp comment '创建时间',
    update_time datetime    not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (group_id),
    unique key uk_name (name)
) comment '资源分组表';