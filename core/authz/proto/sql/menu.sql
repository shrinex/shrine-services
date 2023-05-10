-- goctl model mysql ddl -src="./core/authz/proto/sql/menu.sql" -dir="./core/authz/proto/model" -c --style=go_zero

create database if not exists authz;
use authz;

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