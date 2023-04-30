package page

func NumPages(total, pageSize int64) int64 {
	if total < 1 || pageSize < 1 {
		return 0
	}

	mod := total % pageSize
	if mod == 0 {
		return total / pageSize
	}

	return total/pageSize + 1
}
