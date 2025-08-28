package web

import (
	"encoding/json"
	"net/http"
)

type Map map[string]any

func Success(w http.ResponseWriter, m Map) {
	// 设置响应头为 JSON 格式
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 将 map 序列化为 JSON 并写入响应
	json.NewEncoder(w).Encode(m)
}

func Failure(w http.ResponseWriter, m Map) {
	// 设置响应头为 JSON 格式
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	// 将 map 序列化为 JSON 并写入响应
	json.NewEncoder(w).Encode(m)
}
