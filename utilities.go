package main

import "sort"

func stringSliceContainsString(strSlice []string, str string) bool {
	sort.Strings(strSlice)
	i := sort.SearchStrings(strSlice, str)
	return i < len(strSlice) && strSlice[i] == str
}
