package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/exec"
	"reflect"
	"strings"
)

type Map map[string]any

// bind 解析body参数
func Bind(r *http.Request, v any) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(v)
	if err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	val := reflect.ValueOf(v).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		fieldType := typ.Field(i)
		jsonTag := fieldType.Tag.Get("json")

		// 如果JSON标签包含omitempty，则允许为空
		if strings.Contains(jsonTag, "omitempty") {
			continue
		}

		// 判断字段值是否为零值
		if fieldValue.IsZero() {
			return fmt.Errorf("请求体参数[%s]不能为空", fieldType.Name)
		}
	}

	return nil
}

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

func RunCmd(command string) error {
	cmd := exec.Command("bash", "-c", command)
	slog.Info("执行命令", slog.String("Cmd", command))
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("执行命令报错", slog.String("Err", err.Error()))
		return err
	}
	slog.Info("执行命令结果", slog.String("Out", string(output)))
	return nil
}
