//+build darwin linux unix

package util

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
)

var stdErrFileHandler *os.File
// go的运行时错误panic默认写到标准错误中, 由于无法持久化容器重启后无法排查引起panic的原因
// 使用此方法将标准错误重写到持久化的日志文件里
func RewritePanicsToFile(topic string) error {
	errFile := "/home/golanger/log/" + topic + "_stdErr.log"
	errFileHandler, err := os.OpenFile(errFile, os.O_RDWR|os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	stdErrFileHandler = errFileHandler //把文件句柄保存到全局变量，避免被GC回收

	if err = syscall.Dup2(int(errFileHandler.Fd()), int(os.Stderr.Fd())); err != nil {
		return err
	}
	// GC前关闭文件描述符
	runtime.SetFinalizer(stdErrFileHandler, func(fd *os.File) {
		fd.Close()
	})

	return nil
}

