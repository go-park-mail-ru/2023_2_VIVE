package utils

func Contains(elemToCheckFor int, elements []int) bool {
	for _, elem := range elements {
		if elem == elemToCheckFor {
			return true
		}
	}
	return false
}

// Returns difference between two slices of ints.
// For example: {1, 2, 3} - {1, 2} = {3}
func Difference(slice1, slice2 []int) []int {
	m := make(map[int]bool)

	for _, item := range slice2 {
		m[item] = true
	}

	res := []int{}

	for _, item := range slice1 {
		if _, ok := m[item]; !ok {
			res = append(res, item)
		}
	}

	return res
}
