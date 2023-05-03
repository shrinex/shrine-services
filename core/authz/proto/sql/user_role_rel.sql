-- goctl model mysql ddl -src="./proto/sql/user_role_rel.sql" -dir="./proto/model" --style=go_zero

create database if not exists authz;
use authz;

drop table if exists user_role_rel;
create table user_role_rel
(
    rel_id      bigint   not null comment '关联ID',
    user_id     bigint   not null comment '用户ID',
    role_id     bigint   not null comment '角色ID',
    create_time datetime not null default current_timestamp comment '创建时间',
    update_time datetime not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (rel_id),
    unique key uk_userid_roleid (user_id, role_id)
) comment '用户-角色关联表';