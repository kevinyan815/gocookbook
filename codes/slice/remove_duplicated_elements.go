package main

import (
    "fmt"
)

func main() {
    s := []string{"hello", "world", "hello", "golang", "hello", "ruby", "php", "java"}

    fmt.Println(removeDuplicateElement(s)) //output: hello world golang ruby php java
}

func removeDuplicateElement(languages []string) []string {
    result := make([]string, 0, len(languages))
    temp := map[string]struct{}{}
    for _, item := range languages {
        if _, ok := temp[item]; !ok {
            temp[item] = struct{}{}
            result = append(result, item)
        }
    }
    return result
}

// output: ["hello", "world", "golang", "ruby", "php", "java"]
