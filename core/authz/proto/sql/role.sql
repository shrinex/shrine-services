-- goctl model mysql ddl -src="./core/authz/proto/sql/role.sql" -dir="./core/authz/proto/model" -c --style=go_zero

create database if not exists authz;
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


