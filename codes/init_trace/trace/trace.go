package trace

import "fmt"

func Trace(t string, v int) int {
	fmt.Println(t, ":", v)
	return v
}
