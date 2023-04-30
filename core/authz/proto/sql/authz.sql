-- goctl model mysql ddl -src="./proto/sql/*.sql" -dir="./proto/model" -c =go_zero

drop database if exists authz;
create database authz;
use authz;

drop table if exists role;
create table role
(
    role_id     bigint      not null comment '角色ID',
    name        varchar(12) not null comment '角色名称',
    remark      varchar(36) not null comment '备注',
    creator_id  bigint      not null comment '创建者ID',
    shop_id     bigint      not null comment '所属店铺ID',
    sys_type    tinyint     not null comment '系统类型: 1-平台端,2-商家端,3-普通用户',
    create_time datetime    not null default current_timestamp comment '创建时间',
    update_time datetime    not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (role_id),
    unique key uk_name (name)
) comment '角色表';

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

drop table if exists menu;
create table menu
(
    menu_id     bigint       not null comment '菜单ID',
    name        varchar(12)  not null comment '菜单名称',
    icon        varchar(128) not null comment '图标',
    parent_id   bigint       not null comment '父级菜单ID,0表示一级菜单',
    level       tinyint      not null default 0 comment '层级[1-3]',
    path        varchar(128) not null comment '路径,从一级菜单ID到当前菜单ID',
    sys_type    tinyint      not null comment '系统类型: 1-平台端,2-商家端,3-普通用户',
    weight      int          not null default 0 comment '权重',
    create_time datetime     not null default current_timestamp comment '创建时间',
    update_time datetime     not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (menu_id),
    unique key uk_name (name),
    unique key uk_path (path)
) comment '菜单表';

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


