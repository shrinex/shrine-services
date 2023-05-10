-- goctl model mysql ddl -src="./core/authz/proto/sql/role_menu_rel.sql" -dir="./core/authz/proto/model" --style=go_zero

create database if not exists authz;
use authz;

drop table if exists role_menu_rel;
create table role_menu_rel
(
    rel_id      bigint   not null comment '关联ID',
    role_id     bigint   not null comment '角色ID',
    menu_id     bigint   not null comment '菜单ID',
    create_time datetime not null default current_timestamp comment '创建时间',
    update_time datetime not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (rel_id),
    unique key uk_roleid_menuid (role_id, menu_id)
) comment '角色-菜单关联表';