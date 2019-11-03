package common

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"os"
)

func dumpDb(dbSlice []string, outDir string, engine *xorm.Engine)  {
	for _, db := range dbSlice {
		sql := fmt.Sprintf("use %s;", db)
		_, err := engine.QueryString(sql)
		if err != nil {
			// 说明访问数据库异常，可能不存在
			fmt.Println("数据库", db, "不存在,err is", err)
			os.Exit(0)
		}
		// 说明数据库正常，将数据库信息纪录文件即可
		// 根据数据库名称拼接备份文件
		dumpFile := fmt.Sprintf("%s/%s-create-schema.sql", outDir, db)
		data := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;\r\n", db)
		if writeData(dumpFile, data) {
			fmt.Println("dump schema info of", db, "is ok")
		} else {
			fmt.Println("dump schema info of", db, "is bad")
		}
	}
}
