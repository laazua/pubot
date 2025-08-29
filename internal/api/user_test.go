package api

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"
)

// 封装发送http请求
func req(method, uri string, data []byte) ([]byte, error) {
	var err error
	var req *http.Request
	client := &http.Client{Timeout: time.Duration(time.Second * 30)}
	if data != nil {
		req, err = http.NewRequest(method, "http://127.0.0.1:8088"+uri, bytes.NewBuffer(data))
	} else {
		req, err = http.NewRequest(method, "http://127.0.0.1:8088"+uri, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("resposne code not ok")
	}
	body, err := io.ReadAll(resp.Body)
	return body, err
}

func TestUserCreate(t *testing.T) {
	user := []byte(`{"name": "wangwu", "email": "2222222222@qq.com", "password":"123456abc", "is_admin": false}`)
	body, err := req(http.MethodPost, "/api/user", user)
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("Response Body: %v", string(body))
}

func TestUserDelete(t *testing.T) {
	body, err := req(http.MethodDelete, "/api/user/3", nil)
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("Response Body: %v", string(body))
}

func TestUserUpdate(t *testing.T) {
	user := []byte(`{"id": 2, "name": "wangerxiao", "email": "15132125421@qq.com", "password":"123abc"}`)
	body, err := req(http.MethodPut, "/api/user", user)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(string(body))
}

func TestUserGet(t *testing.T) {
	body, err := req(http.MethodGet, "/api/user/1", nil)
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("Response Body: %v", string(body))
}

func TestUserList(t *testing.T) {
	body, err := req(http.MethodGet, "/api/users", nil)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(string(body))
}
