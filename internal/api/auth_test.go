package api

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestUserLogin(t *testing.T) {
	user := []byte(`{"name": "admin", "password": "123456"}`)
	client := http.Client{Timeout: time.Duration(time.Second * 30)}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8088/api/login", bytes.NewBuffer(user))
	if err != nil {
		t.Error(err.Error())
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Error(errors.New("server is error"))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(errors.New("read response body errror"))
	}
	t.Log(string(body))
}
