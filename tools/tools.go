package tools

func Min64(a uint64, b uint64) uint64 {
	if a < b {
		return a
	} else if a == b {
		return a
	} else {
		return b
	}
}

func Max64(a uint64, b uint64) uint64 {
	if a > b {
		return a
	} else if a == b {
		return a
	} else {
		return b
	}
}

func SliceIndex(limit int, predicate func(i int) bool) (pos int, ok bool) {
	pos = -1
	ok = false
	for i := 0; i < limit; i++ {
		if predicate(i) {
			pos = i
			ok = true
			return
		}
	}
	return
}
