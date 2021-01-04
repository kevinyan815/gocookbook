import (
    "fmt"
    "reflect"
    "unsafe"
)
 
func main() {
 
    v1 := uint(12)
    v2 := int(13)
 
    fmt.Println(reflect.TypeOf(v1)) //uint
    fmt.Println(reflect.TypeOf(v2)) //int
 
    fmt.Println(reflect.TypeOf(&v1)) //*uint
    fmt.Println(reflect.TypeOf(&v2)) //*int
 
    p := &v1
    p = (*uint)(unsafe.Pointer(&v2)) //使用unsafe.Pointer进行类型的转换
 
    fmt.Println(reflect.TypeOf(p)) // *unit
    fmt.Println(*p) //13
}
