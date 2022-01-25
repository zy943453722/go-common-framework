package util

import "testing"

func TestGetFirstDayOfMonth(t *testing.T) {
	res, err := GetFirstDayOfMonth("2021-12", "2006-01")
	if err != nil {
		t.Error(err)
	}
	t.Logf("结果为:%v", res)
}

func TestGetLastDayOfMonth(t *testing.T) {
	res, err := GetLastDayOfMonth("2021-12", "2006-01")
	if err != nil {
		t.Error(err)
	}
	t.Logf("结果为:%v", res)
}

func TestParseTimeToDate(t *testing.T) {
	res, err := ParseTimeToDate("2020-09-19 00:00:00", "2006-01-02 15:04:05", "2006-01-02")
	if err != nil {
		t.Error(err)
	}
	t.Logf("结果为:%v", res)
}

func TestGetDateSubDays(t *testing.T) {
	dateStart, err := GetFirstDayOfMonth("2021-12", "2006-01")
	if err != nil {
		t.Error(err)
	}
	dateEnd, err := GetLastDayOfMonth("2021-12", "2006-01")
	if err != nil {
		t.Error(err)
	}
	res := GetDateSubDays(dateStart, dateEnd)
	t.Logf("结果为:%v", res)
}

func TestGetMonthDays(t *testing.T) {
	days, err := GetMonthDays("2021-11", "2006-01")
	if err != nil {
		t.Error(err)
	}
	t.Logf("结果为:%d", days)
}

func TestGetMonthCount(t *testing.T) {
	months, err := GetMonthCount("2021-12", "2023-11", "2006-01")
	if err != nil {
		t.Error(err)
	}
	t.Logf("结果为:%d", months)
}

func TestGetMonthSub(t *testing.T) {
	months, err := GetMonthSub("2021-12", "2023-11", "2006-01")
	if err != nil {
		t.Error(err)
	}
	t.Logf("结果为:%v", months)
}
