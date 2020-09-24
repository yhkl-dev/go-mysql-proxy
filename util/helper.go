package util

import (
	"fmt"
	"strings"

	mapset "github.com/deckarep/golang-set"
)

func Contains(slice []interface{}, item interface{}) bool {
	set := mapset.NewSetFromSlice(slice)
	return set.Contains(item)
}

func IntersectSlice(s1 []interface{}, s2 []interface{}) []interface{} {
	if s1 == nil && s2 == nil {
		return nil
	}
	if s1 == nil {
		return s2
	}

	if s2 == nil {
		return s1
	}
	set1 := mapset.NewSetFromSlice(s1)
	set2 := mapset.NewSetFromSlice(s2)

	return set1.Intersect(set2).ToSlice()
}

func implode(list []interface{}) {
	strings.Replace(strings.Trim(fmt.Sprint(list), "[]"), " ", ",", -1)
}
