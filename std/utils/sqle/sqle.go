package sqle

import "github.com/go-sql-driver/mysql"

const (
	DuplicateEntry uint16 = 1062
)

func Is(err error, code uint16) bool {
	if err == nil {
		return false
	}

	myerr, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}

	return myerr.Number == code
}
