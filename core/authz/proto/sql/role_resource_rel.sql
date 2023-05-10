-- goctl model mysql ddl -src="./core/authz/proto/sql/role_resource_rel.sql" -dir="./core/authz/proto/model" --style=go_zero

create database if not exists authz;
use authz;

drop table if exists role_resource_rel;
create table role_resource_rel
(
    rel_id      bigint   not null comment '关联ID',
    role_id     bigint   not null comment '角色ID',
    resource_id bigint   not null comment '资源ID',
    create_time datetime not null default current_timestamp comment '创建时间',
    update_time datetime not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (rel_id),
    unique key uk_roleid_resourceid (role_id, resource_id)
) comment '角色-资源关联表';