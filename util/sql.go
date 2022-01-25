package util

import "strings"

func GetInSqlStr(str string) string {
	if str == "" {
		return ""
	}
	arr := strings.Split(str, ",")
	temp := ""
	for _, v := range arr {
		temp += `'` + v + `',`
	}
	temp = strings.Replace(temp, `''`, `'`, -1)
	return strings.TrimSuffix(temp, ",")
}
