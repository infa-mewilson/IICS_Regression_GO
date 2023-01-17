package utils

import "sort"

func SortingMap(m map[string]int) []string {

	//using a sorted slice of keys to return a map[string]int in key order.
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
