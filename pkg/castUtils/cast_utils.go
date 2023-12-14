package castUtils

func IntSliceToInt64Slice(slice []int) []int64 {
	res := make([]int64, len(slice))
	for i, item := range slice {
		res[i] = int64(item)
	}
	return res
}

func Int64SliceToIntSlice(slice []int64) []int {
	res := make([]int, len(slice))
	for i, item := range slice {
		res[i] = int(item)
	}
	return res
}

func StringToAnySlice(slice []string) []interface{} {
	res := make([]interface{}, len(slice))
	for i, item := range slice {
		res[i] = item
	}
	return res
}
