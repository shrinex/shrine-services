package verify

import "github.com/go-sql-driver/mysql"

func Duplicated(err error) bool {
	if err == nil {
		return false
	}

	myerr, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}

	return myerr.Number == 1062
}
