package main

import (
	"example.com/init_trace/pkg1"
	"example.com/init_trace/pkg2"

	"example.com/init_trace/trace"
	"fmt"
)

func init() {
	fmt.Println("init1 func in main")
}

func init() {
	fmt.Println("init2 func in main")
}

var M_v1 = trace.Trace("init M_v1", pkg1.P1_v2 + 10)
var M_v2 = trace.Trace("init M_v2", pkg2.P2_v2 + 10)


func main() {
	fmt.Println("main func in main")
}
