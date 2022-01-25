package util

import (
	"encoding/json"
	"testing"
)

type UserInfo struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"fullname"`
}

func TestAesCBCEncrypt(t *testing.T) {
	user := &UserInfo{
		UserName: "xxxxx",
		Email:    "xxxx@qq.com",
		FullName: "xxxx",
	}
	key := "xxxxxxxxxxx"
	raw, _ := json.Marshal(user)
	res, _ := AesCBCEncrypt(raw, []byte(key))
	t.Logf("结果为:%s", Base64Encrypt(res))
}

func TestAesCBCDecrypt(t *testing.T) {
	realData := "xxxxxx"
	key := "xxxxxxxxx"
	data, err := Base64Decrypt(realData)
	if err != nil {
		t.Error(err)
	}
	res, err := AesCBCDecrypt(data, []byte(key))
	if err != nil {
		t.Error(err)
	}
	user := new(UserInfo)
	if err = json.Unmarshal(res, user); err != nil {
		t.Error(err)
	}
	t.Logf("结果为:%v", user)
}

func TestGetMd5(t *testing.T) {
	pin := "123"
	salt := "xxxxx"
	res := GetMd5(pin, salt)
	t.Logf("结果为:%s", res)
}

func TestGetSha512(t *testing.T) {
	str := "123aaaa"
	res := GetSha512(str)
	t.Logf("结果为:%s", res)
}
