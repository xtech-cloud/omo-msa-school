package tool

func Equal(source []string, dest []string) bool {
	if source == nil || dest == nil {
		return false
	}
	if len(source) != len(dest) {
		return false
	}
	for _, s := range source {
		for i := 0;i < len(dest);i += 1 {
			if s != dest[i] {
				return false
			}
		}
	}
	return true
}
