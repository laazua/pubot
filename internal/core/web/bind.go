package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

func Bind(r *http.Request, v any) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	defer r.Body.Close()

	// 解码 JSON 数据
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(v)
	if err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	val := reflect.ValueOf(v).Elem()
	// 遍历结构体字段
	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		// 判断字段值是否为零值
		if fieldValue.IsZero() {
			return fmt.Errorf("请求体参数[%v]不能为空,当前值[%v]", val.Type().Field(i).Name, fieldValue)
		}
	}

	return nil
}
