package utils

func GetSwitchFromOptions(Options map[string]bool, key string) (result bool) {
	if flag, ok := Options[key]; !ok || flag {
		return true
	}
	return false
}

func SetSwitchFromOptions(options map[string]bool, key string, value bool) {
	if options == nil {
		options = make(map[string]bool, 5)
	}
	options[key] = value
}
