package annotator

func packageName(name string) string {
	for i := len(name) - 1; i >= 0 && name[i] != '/'; i-- {
		if name[i] == '.' {
			return name[:i]
		}
	}
	return ""
}
