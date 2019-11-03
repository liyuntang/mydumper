package common

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"os"
	"strings"
)

func dumpTabInfo(outDir, tableInfo string, engine *xorm.Engine)  {
	sql := fmt.Sprintf("show create table %s;", tableInfo)
	res, err := engine.QueryString(sql)
	if err != nil {
		fmt.Println("get table info of", tableInfo, "is bad, err is", err)
		os.Exit(0)
	}
	// 拼接备份文件名称
	dumpFile := fmt.Sprintf("%s/%s-create-table.sql", outDir, strings.Replace(tableInfo, ".", "-",1))
	for _, dict := range res {
		data := fmt.Sprintf("%s;\r\n", dict["Create Table"])
		writeData(dumpFile, data)
	}
}
