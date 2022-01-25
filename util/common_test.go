package util

import "testing"

func TestRemoveDuplication(t *testing.T) {
	arr := [][]string{
		{"SN3334343", "1", "2", "3"},
		{"SNerehrth", "1", "2", "3"},
		{"SN3334343", "4", "5", "6"},
	}
	res := RemoveDuplication(arr, 0)
	t.Logf("%v", res)
}

func TestImplode(t *testing.T) {
	arr := []string{"a", "b", "c"}
	res := Implode(",", arr)
	t.Logf("res:%s", res)
}
