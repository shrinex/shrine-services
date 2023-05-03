package rdb

import (
	"github.com/go-sql-driver/mysql"
	"time"
)

type MySQLConf struct {
	User   string // Username
	Passwd string // Password (requires User)
	Addr   string // Network address (requires Net)
	DBName string // Database name
}

func (ms *MySQLConf) FormatDSN() string {
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = ms.Addr
	cfg.User = ms.User
	cfg.Loc = time.Local
	cfg.Passwd = ms.Passwd
	cfg.DBName = ms.DBName
	cfg.ParseTime = true
	cfg.MultiStatements = true
	cfg.AllowOldPasswords = true
	cfg.Collation = "utf8mb4_unicode_ci"
	return cfg.FormatDSN()
}
