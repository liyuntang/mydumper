package common

import (
	"fmt"
	"io"
	"os"
)

func writeData(dataFile, data string) (isok bool) {
	// 判断dataFile文件是否存在，根据存在的状态使用不用的flag
	var flag int
	_, err := os.Stat(dataFile)
	if err != nil {
		flag = os.O_CREATE|os.O_WRONLY
	} else {
		flag = os.O_APPEND|os.O_WRONLY
	}
	file, err := os.OpenFile(dataFile, flag, 0644)
	defer file.Close()
	if err != nil {
		fmt.Println("create dumpFile", dataFile, "is bad, err is", err)
		return false
	}
	_, err1 := io.WriteString(file, data)
	if err1 != nil {
		fmt.Println("write data to dumpFile is bad, err is", err1)
		return false
	}
	return true
}
