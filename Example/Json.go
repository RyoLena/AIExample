package main

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
)

// parseRequestJSON 会用请求 主体中 JSON 编码值的字段填充目标。
// 它希望请求的 Content-Type 标头设置为 JSON，
// 请求体的 JSON 编码值符合目标的底层类型。
func parseRequestJSON(r *http.Request, target any) error {
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		fmt.Println("解析的格式 不对")
		return err
	}
	if mediaType != "application/json" {
		fmt.Println("意外的内容格式")
		return fmt.Errorf("expected Content-Type application/json; got %s", contentType)
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	return dec.Decode(target)
}

// renderResponseJSON 将 res 编码为 JSON 并写入 w。
func renderResponseJSON(w http.ResponseWriter, res any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
