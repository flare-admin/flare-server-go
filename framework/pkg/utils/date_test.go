package utils

import "testing"

func Test_GetStartAndEndOfDay(t *testing.T) {
	day, t2 := GetStartAndEndOfDay()
	t.Log(day, t2)
}
