package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

func RunCmd(command, workDir string) error {
	cmd := exec.Command("bash", "-c", command)
	if workDir != "" {
		cmd.Dir = workDir
	}

	slog.Info("执行命令", slog.String("Cmd", command), slog.String("Dir", cmd.Dir))

	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("执行命令报错", slog.String("Err", err.Error()))
		return err
	}
	slog.Info("执行命令结果", slog.String("Out", string(output)))
	return nil
}

// RunCommands 支持 cd 持久化
func RunCommands(cmds []string) error {
	// 当前目录（初始为程序启动目录）
	currentDir, _ := os.Getwd()

	for _, raw := range cmds {
		c := strings.TrimSpace(raw)
		if c == "" {
			continue
		}

		// 处理 cd
		if strings.HasPrefix(c, "cd ") {
			dir := strings.TrimSpace(strings.TrimPrefix(c, "cd "))
			// 转换成绝对路径
			if !filepath.IsAbs(dir) {
				dir = filepath.Join(currentDir, dir)
			}
			if _, err := os.Stat(dir); err != nil {
				return errors.New("目录不存在: " + dir)
			}
			currentDir = dir
			slog.Info("切换目录", slog.String("Dir", currentDir))
			continue
		}

		// 普通命令
		if err := RunCmd(c, currentDir); err != nil {
			return err
		}
	}
	return nil
}

func ChWorkSpace(path string) error {
	err := os.Chdir(path)
	if err != nil {
		slog.Error("切换工作目录失败,检查是否配置")
		return err
	}
	return nil
}
