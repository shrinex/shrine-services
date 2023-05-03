-- goctl model mysql ddl -src="./proto/sql/resource.sql" -dir="./proto/model" -c --style=go_zero

create database if not exists authz;
use authz;

drop table if exists resource;
create table resource
(
    resource_id bigint       not null comment '资源ID',
    group_id    bigint       not null comment '所属分组ID',
    name        varchar(12)  not null comment '资源名称',
    method      varchar(7)   not null comment '请求方法',
    pattern     varchar(128) not null comment '资源路径（ant style）',
    sys_type    tinyint      not null comment '系统类型: 1-平台端,2-商家端,3-普通用户',
    create_time datetime     not null default current_timestamp comment '创建时间',
    update_time datetime     not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (resource_id),
    key idx_groupid (group_id),
    unique key uk_name (name),
    unique key uk_systype_method_pattern (sys_type, method, pattern)
) comment '资源表';