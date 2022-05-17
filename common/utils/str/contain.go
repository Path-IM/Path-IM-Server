package strUtils

func IsContain(target string, List []string) bool {
	for _, element := range List {

		if target == element {
			return true
		}
	}
	return false
}
