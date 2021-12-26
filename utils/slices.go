package utils

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
