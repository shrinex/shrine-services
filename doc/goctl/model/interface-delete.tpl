TxDelete(ctx context.Context, tx sqlx.Session, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error
Delete(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error