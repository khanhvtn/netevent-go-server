package utilities

func Contains(element string, list []string) bool {
	for _, ele := range list {
		if ele == element {
			return true
		}
	}
	return false
}
