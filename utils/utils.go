package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

// PrettyJSON 美化打印
func PrettyJSON(v interface{}) {
	// 使用 json.MarshalIndent 进行格式化和美化打印
	prettyJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshalling failed: %s", err)
	}
	// 打印格式化后的 JSON 字符串
	fmt.Println(string(prettyJSON))
}

func ToPrettyJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}
