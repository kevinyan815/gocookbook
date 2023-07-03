package slice_util

func StringInSlice(search string, arr []string) bool {
	for _, v := range arr {
		if v == search {
			return true
		}
	}
	return false
}

func FilterStringSlice(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
