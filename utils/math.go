package utils

import "sort"

func MedianSlice(sl []uint64) uint64 {
	l := len(sl)

	if l == 0 {
		return 0
	}

	if l == 1 {
		return sl[0]
	}

	sort.Slice(sl, func(i, j int) bool { return sl[i] < sl[j] })

	n := l / 2
	if l%2 == 0 {
		return (sl[n-1] + sl[n]) / 2
	} else {
		return sl[n]
	}
}
