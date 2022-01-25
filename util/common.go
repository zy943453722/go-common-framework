package util

import (
	"strings"

	"github.com/gofrs/uuid"
)

func GenUuid() string {
	return uuid.Must(uuid.NewV4()).String()
}

func GenResourceUuid(prefix string) string {
	u := GenUuid()
	return prefix + "_" + strings.Replace(u, "-", "", -1)[:8]
}

func GetMapKey(data map[string]string) []string {
	s := make([]string, 0)
	for key := range data {
		s = append(s, key)
	}
	return s
}

func GetMapValue(data map[string]string) []string {
	s := make([]string, 0)
	for _, value := range data {
		s = append(s, value)
	}
	return s
}

func InArray(value string, arr []string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

//RemoveDuplication excel的二维数组去重
func RemoveDuplication(arr [][]string, offset int) [][]string {
	set := make(map[string]int, len(arr))
	j := 0
	for _, v := range arr {
		v[offset] = strings.TrimSpace(v[offset])
		if value, ok := set[v[offset]]; ok {
			arr[value] = v
			continue
		}
		set[v[offset]] = j
		arr[j] = v
		j++
	}
	return arr[:j]
}

//Implode 切片根据指定字符转字符串
func Implode(glue string, data []string) string {
	if len(data) <= 0 {
		return ""
	}
	res := ""
	for _, v := range data {
		res += v + glue
	}
	return strings.TrimRight(res, glue)
}

//ArrayFlip map的key-value互换
func ArrayFlip(m map[string]string) map[string]string {
	n := make(map[string]string)
	for i, v := range m {
		n[v] = i
	}
	return n
}
