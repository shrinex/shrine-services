package authx

import "strconv"

func parseInt(in string) int64 {
	ret, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		panic(err)
	}

	return ret
}

func formatInt(in int64) string {
	return strconv.FormatInt(in, 10)
}
