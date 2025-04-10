package common

import (
	"fmt"
	"strings"
)

func CheckIfSliceContainStr(a string, b []string) bool {
	for _, v := range b {
		if a == v {
			return true
		}
	}
	return false
}

func CheckStringArrOverlap(a []string, b []string) bool {
	for _, va := range a {
		for _, vb := range b {
			if va == vb {
				return true
			}
		}
	}
	return false
}

func ConvertNumArrToString(a []uint) string {
	if len(a) == 0 {
		return ""
	}
	idsStr := ""
	for _, id := range a {
		idsStr += fmt.Sprintf("%v,", id)

	}
	return idsStr[:len(idsStr)-1]
}

func ParseSortString(sortStr string) string {
	sort := strings.Split(sortStr, ".")
	if len(sort) == 1 {
		return sort[0] + " ASC"
	}
	return sort[0] + " " + sort[1]
}