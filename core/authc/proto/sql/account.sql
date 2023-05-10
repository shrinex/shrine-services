-- goctl model mysql ddl -src="./core/authc/proto/sql/account.sql" -dir="./core/authc/proto/model" -c --style=go_zero

create database if not exists authc;
use authc;

drop table if exists account;
create table account
(
    account_id  bigint      not null comment '账号ID',
    user_id     bigint      not null comment '关联的用户ID',
    username    varchar(32) not null comment '用户名',
    password    varchar(64) not null comment '密码',
    sys_type    tinyint     not null comment '系统类型: 1-平台端,2-商家端,3-普通用户',
    shop_id     bigint      not null comment '所属店铺,<=0表示未关联店铺',
    is_admin    tinyint     not null comment '是否是管理员',
    enabled     tinyint     not null comment '状态：0-禁用,1-启用',
    create_time datetime    not null default current_timestamp comment '创建时间',
    update_time datetime    not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (account_id),
    unique key uk_systype_userid (sys_type, user_id),
    unique key uk_systype_username (sys_type, username)
) comment '账号表';


