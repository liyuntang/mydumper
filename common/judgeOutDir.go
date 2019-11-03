package common

import (
	"fmt"
	"io/ioutil"
	"os"
)

func judgeOutDir(dir string) {
	// 查看目录是否存在
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		// 说明dir不存在，需要创建该dir
		err1 := os.MkdirAll(dir, 0755)
		if err1 != nil {
			fmt.Println("create dir is bad, err1 is", err1)
			os.Exit(0)
		}
	}
	// 说明dir存在，此时需要判断该dir是否为空
	dirSlice, err2 := ioutil.ReadDir(dir)
	if err2 != nil {
		fmt.Println("get dir status is bad, err2 is", err2)
		os.Exit(0)
	}
	if len(dirSlice) >=1 {
		fmt.Println("sorry dir is not null")
		//os.Exit(0)
	}
}
