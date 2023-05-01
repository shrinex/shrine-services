-- goctl model mysql ddl -src="./proto/sql/*.sql" -dir="./proto/model" -c --style=go_zero

create database if not exists merchant;
use merchant;

create table shop
(
    shop_id     bigint       not null comment '店铺ID',
    name        varchar(24)  not null comment '店铺名称',
    intro       varchar(64)  not null comment '店铺简介',
    logo        varchar(128) not null comment '店铺logo',
    status      tinyint      not null comment '店铺状态:0-歇业,1-营业中,2-已下架',
    type        tinyint      not null comment '店铺类型:1-自营店,2-普通店',
    create_time datetime     not null default current_timestamp comment '创建时间',
    update_time datetime     not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (shop_id),
    unique key uk_name (name)
) comment ='店铺表';

