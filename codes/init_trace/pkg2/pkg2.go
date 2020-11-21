package pkg2


import (
	"example.com/init_trace/trace"
	"fmt"
)

var P2_v1 = trace.Trace("init P2_v1", 20)
var P2_v2 = trace.Trace("init P2_v2", 30)

func init() {
	fmt.Println("init func in pkg2")
}
