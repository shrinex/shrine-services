-- goctl model mysql ddl -src="./proto/sql/brand_category_rel.sql" -dir="./proto/model" --style=go_zero

create database if not exists product;
use product;

drop table if exists brand_category_rel;
create table brand_category_rel
(
    rel_id      bigint   not null comment '关联ID',
    category_id bigint   not null comment '分类ID',
    brand_id    bigint   not null comment '品牌ID',
    create_time datetime not null default current_timestamp comment '创建时间',
    update_time datetime not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (rel_id),
    key idx_category_id (category_id),
    unique key uk_brandid_catid (brand_id, category_id)
) comment ='品牌与分类关联表';
