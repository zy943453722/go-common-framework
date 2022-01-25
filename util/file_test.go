package util

import "testing"

func TestGetCurDir(t *testing.T) {
	if str, err := GetBaseDir(); err != nil {
		t.Error(err)
	} else {
		t.Log(str)
	}
}
