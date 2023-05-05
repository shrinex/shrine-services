TxInsert(ctx context.Context, tx sqlx.Session, data *{{.upperStartCamelObject}}) (sql.Result,error)
Insert(ctx context.Context, data *{{.upperStartCamelObject}}) (sql.Result,error)