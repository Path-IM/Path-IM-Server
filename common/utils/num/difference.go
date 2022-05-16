package numUtils

//Get the intersection of two slices
func IntersectUInt32(slice1, slice2 []uint32) []uint32 {
	m := make(map[uint32]bool)
	n := make([]uint32, 0)
	for _, v := range slice1 {
		m[v] = true
	}
	for _, v := range slice2 {
		flag, _ := m[v]
		if flag {
			n = append(n, v)
		}
	}
	return n
}

//Get the diff of two slices
func DifferenceUInt32(slice1, slice2 []uint32) []uint32 {
	m := make(map[uint32]bool)
	n := make([]uint32, 0)
	inter := IntersectUInt32(slice1, slice2)
	for _, v := range inter {
		m[v] = true
	}
	for _, v := range slice1 {
		if !m[v] {
			n = append(n, v)
		}
	}

	for _, v := range slice2 {
		if !m[v] {
			n = append(n, v)
		}
	}
	return n
}
