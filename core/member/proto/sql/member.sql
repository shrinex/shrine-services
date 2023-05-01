-- goctl model mysql ddl -src="./proto/sql/*.sql" -dir="./proto/model" -c --style=go_zero

create database if not exists authc;
use authc;

drop table if exists user;
create table user
(
    user_id     bigint       not null comment '用户ID',
    shop_id     bigint       not null comment '所属店铺,<=0表示未关联店铺',
    sys_type    tinyint      not null comment '系统类型: 1-平台端,2-商家端,3-普通用户',
    nickname    varchar(12)  not null comment '昵称',
    avatar      varchar(128) not null comment '头像',
    intro       varchar(36)  not null default '' comment '一句话介绍自己',
    active      tinyint      not null comment '是否已激活,0-未激活,1-已激活',
    enabled     tinyint      not null comment '状态：0-禁用,1-启用',
    create_time datetime     not null default current_timestamp comment '创建时间',
    update_time datetime     not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (user_id),
    unique key uk_systype_nickname (sys_type, nickname)
) comment '用户表';