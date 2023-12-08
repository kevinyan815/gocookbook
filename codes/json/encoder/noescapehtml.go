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
		Url: "https://yourexample.com/add?type=1&sku=iphone",
	}
	enc, _ := MarshalWithNoEscape(&data)
	fmt.Println(string(enc))
}

func MarshalWithNoEscape(i interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(i)
	return bytes.TrimRight(buffer.Bytes(), "\n"), err
}
