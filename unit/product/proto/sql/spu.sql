-- goctl model mysql ddl -src="./proto/sql/attr.sql" -dir="./proto/model" -c --style=go_zero

create database if not exists product;
use product;

drop table if exists spu;
create table spu
(
    spu_id      bigint      not null comment 'ID',
    shop_id     bigint      not null comment '店铺ID',
    brand_id    bigint               default null comment '品牌ID',
    category_id bigint      not null comment '三级分类ID',
    name        varchar(64) not null comment '商品名',
    status      tinyint     not null comment '状态:0-待审核,1-已上架,2-已失效',
    weight      int         not null comment '排序',
    create_time datetime    not null default current_timestamp comment '创建时间',
    update_time datetime    not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (spu_id)
) comment 'standard product unit';