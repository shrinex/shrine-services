package dtmx

import (
	"database/sql"
	"github.com/zeromicro/go-zero/core/syncx"
	"io"
	"shrine/std/conf/rdb"
	"time"
)

const (
	maxIdleConns = 12
	maxOpenConns = 12
	maxLifetime  = time.Minute
)

var connManager = syncx.NewResourceManager()

func getMySql(cfg rdb.MySQLConf) (*sql.DB, error) {
	server := cfg.FormatDSN()
	val, err := connManager.GetResource(server, func() (io.Closer, error) {
		return newMySql(server)
	})
	if err != nil {
		return nil, err
	}

	return val.(*sql.DB), nil
}

func newMySql(datasource string) (*sql.DB, error) {
	conn, err := sql.Open("mysql", datasource)
	if err != nil {
		return nil, err
	}

	// we need to do this until the issue https://github.com/golang/go/issues/9851 get fixed
	// discussed here https://github.com/go-sql-driver/mysql/issues/257
	// if the discussed SetMaxIdleTimeout methods added, we'll change this behavior
	// 8 means we can't have more than 8 goroutines to concurrently access the same database.
	conn.SetMaxIdleConns(maxIdleConns)
	conn.SetMaxOpenConns(maxOpenConns)
	conn.SetConnMaxLifetime(maxLifetime)

	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return conn, nil
}
