package strUtils

//Get the intersection of two slices
func IntersectString(slice1, slice2 []string) []string {
	m := make(map[string]bool)
	n := make([]string, 0)
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
func DifferenceString(slice1, slice2 []string) []string {
	m := make(map[string]bool)
	n := make([]string, 0)
	inter := IntersectString(slice1, slice2)
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
