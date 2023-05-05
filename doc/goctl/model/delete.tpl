func (m *default{{.upperStartCamelObject}}Model) TxDelete(ctx context.Context, tx sqlx.Session, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	query := fmt.Sprintf("delete from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table)
    _, err := tx.ExecCtx(ctx, query, {{.lowerStartCamelPrimaryKey}})
    if err != nil {
    	return err
    }

	{{if .withCache}}
	    {{if .containsIndexCache}}
	data, err := m.TxFindOne(ctx, tx, {{.lowerStartCamelPrimaryKey}})
	if err != nil {
		return err
	}
        {{end}}
    {{.keys}}
    return m.DelCacheCtx(ctx, {{.keyValues}}){{else}}return err{{end}}
}

func (m *default{{.upperStartCamelObject}}Model) Delete(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err:=m.FindOne(ctx, {{.lowerStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}

{{end}}	{{.keys}}
    _, err {{if .containsIndexCache}}={{else}}:={{end}} m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table)
		return conn.ExecCtx(ctx, query, {{.lowerStartCamelPrimaryKey}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("delete from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table)
		_,err:=m.conn.ExecCtx(ctx, query, {{.lowerStartCamelPrimaryKey}}){{end}}
	return err
}
