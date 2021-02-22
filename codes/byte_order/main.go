package main

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

const INT_SIZE = int(unsafe.Sizeof(0)) //64位操作系统，8 bytes

//判断我们系统中的字节序类型
func systemEdian() {

	var i = 0x01020304
	fmt.Println("&i:",&i)
	bs := (*[INT_SIZE]byte)(unsafe.Pointer(&i))

	if bs[0] == 0x04 {
		fmt.Println("system edian is little endian")
	} else {
		fmt.Println("system edian is big endian")
	}
	fmt.Printf("temp: 0x%x,%v\n",bs[0],&bs[0])
	fmt.Printf("temp: 0x%x,%v\n",bs[1],&bs[1])
	fmt.Printf("temp: 0x%x,%v\n",bs[2],&bs[2])
	fmt.Printf("temp: 0x%x,%v\n",bs[3],&bs[3])

}


func testBigEndian() {

	var testInt int32 = 0x01020304
	fmt.Printf("%d use big endian: \n", testInt)
	testBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(testBytes, uint32(testInt))
	fmt.Println("int32 to bytes:", testBytes)
	fmt.Printf("int32 to bytes: %x \n", testBytes)

	convInt := binary.BigEndian.Uint32(testBytes)
	fmt.Printf("bytes to int32: %d\n\n", convInt)
}

func testLittleEndian() {

	var testInt int32 = 0x01020304
	fmt.Printf("%x use little endian: \n", testInt)
	var testBytes []byte = make([]byte, 4)
	binary.LittleEndian.PutUint32(testBytes, uint32(testInt))
	fmt.Printf("int32 to bytes: %x \n", testBytes)

	convInt := binary.LittleEndian.Uint32(testBytes)
	fmt.Printf("bytes to int32: %d\n\n", convInt)
}

func main() {
	systemEdian()
	fmt.Println("")
	testBigEndian()
	testLittleEndian()
}
