package main
 
import (
    "fmt"
    "unsafe"
)
 
func main() {
 
    var x struct {
        a int
        b int
        c []int
    }
 
    // unsafe.Offsetof 函数的参数必须是一个字段,  比如 x.b,  方法会返回 b 字段相对于 x 起始地址的偏移量, 包括可能的空洞。

    // 指针运算 uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)。
    
    // 和 pb := &x.b 等价
    pb := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))

    *pb = 42
    fmt.Println(x.b) // "42"
}
