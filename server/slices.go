package server

func isSubset(a, b []string) bool {
	found := false
	for _, x := range b {
		for _, y := range a {
			if x == y {
				found = true
				break
			}
		}
		if found == false {
			return false
		}
		found = false
	}
	return true
}
