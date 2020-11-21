package pkg1

import (
	"example.com/init_trace/pkg2"
	"example.com/init_trace/trace"
	"fmt"
)

var P1_v1 = trace.Trace("init P1_v1", pkg2.P2_v1 + 10)
var P1_v2 = trace.Trace("init P1_v2", pkg2.P2_v2 + 10)

func init() {
	fmt.Println("init func in pkg1")
}
