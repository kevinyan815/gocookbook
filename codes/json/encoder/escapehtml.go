package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func main() {
	data := struct {
		Name string
		Url string
	}{
		Name: "Jesper",
		Url: "https://m.qsebao.com/cbb/landpage?scene_name=zhonghui-health-mofang-jf-2023-bundledE3-zrheart-land-page-weiaizf&sku_name=zhonghui-health-mofang-jf-base",
	}
	enc, _ := MarshalWithNoEscape(&data)
	fmt.Println(string(enc))
}

func MarshalWithNoEscape(i interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	//encoder.SetEscapeHTML(false)
	err := encoder.Encode(i)
	return bytes.TrimRight(buffer.Bytes(), "\n"), err
}
