package utils

import "sort"

// SliceIsUniqueUint32 checks if all the values in the slice are unique
func SliceIsUniqueUint32(s *[]uint32) bool {
	uniq := map[uint32]bool{}

	for _, val := range *s {
		if _, ok := uniq[val]; ok {
			return false
		}

		uniq[val] = true
	}

	return true
}

// SliceIsSortedUint32 checks if slice is sorted
func SliceIsSortedUint32(s *[]uint32) bool {
	return sort.SliceIsSorted(*s, func(i, j int) bool {
		return (*s)[i] < (*s)[j]
	})
}

func SliceReverse(s []uint64) []uint64 {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}
