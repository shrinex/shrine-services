func (m *default{{.upperStartCamelObject}}Model) TxInsert(ctx context.Context, tx sqlx.Session, data *{{.upperStartCamelObject}}) (sql.Result,error) {
	query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
  	ret, err := tx.ExecCtx(ctx, query, {{.expressionValues}})
  	if err != nil {
  		return nil, err
  	}
	{{if .withCache}}
    {{.keys}}
  	if err = m.DelCacheCtx(ctx, {{.keyValues}}); err != nil {
  		return nil, err
  	}

  	{{else}}
    {{end}}
	return ret, err
}

func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, data *{{.upperStartCamelObject}}) (sql.Result,error) {
	{{if .withCache}}{{.keys}}
    ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
    ret,err:=m.conn.ExecCtx(ctx, query, {{.expressionValues}}){{end}}
	return ret,err
}
