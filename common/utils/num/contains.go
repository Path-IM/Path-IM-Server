package numUtils

func IsContainInt32(target int32, List []int32) bool {
	for _, element := range List {
		if target == element {
			return true
		}
	}
	return false
}

func IsContainUInt32(target uint32, List []uint32) bool {
	for _, element := range List {
		if target == element {
			return true
		}
	}
	return false
}
