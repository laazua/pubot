package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

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
