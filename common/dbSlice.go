package common

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"os"
	"strings"
)

func getDbSlice(dbs string, engine *xorm.Engine) []string {
	// dbs表示程序接收的要备份的数据库名称，有三种情况：
	// 1、""表示备份所有数据库里的表 2、db1表示只备份db1里的表 3、db1、db2、db3表示备份db1、db2、db3里的数据表
	// 声明一个slice用于存放db信息
	dbSlice := []string{}
	dbLen := len(strings.Split(strings.Trim(dbs, ""), "," ))
	if dbLen == 1 {
		if len(dbs) == 0 {
			// 表示备份所有数据库里的表
			sql := fmt.Sprintf("select distinct(TABLE_SCHEMA) from information_schema.tables where TABLE_SCHEMA not in ('information_schema', 'performance_schema');")
			res, err := engine.QueryString(sql)
			if err != nil {
				// 说明操作数据库失败了
				fmt.Println("获取数据库信息失败，err is", err)
				os.Exit(0)
			}
			for _, dict := range res {
				dbSlice = append(dbSlice, dict["TABLE_SCHEMA"])
			}
		} else {
			// 表示备份指定库的数据表
			// 检查数据库
			checkDb(strings.Trim(dbs, ""), engine)
			dbSlice = append(dbSlice, strings.Trim(dbs, ""))
		}
	} else if dbLen >= 2 {
		// 表示备份某几个库的数据库
		for _, db := range strings.Split(strings.Trim(dbs, ""), "," ){
			// 检查数据库
			checkDb(db, engine)
			dbSlice = append(dbSlice, strings.Trim(db, ""))
		}
	} else {
		// 数据有误，提示用户，退出程序
		fmt.Println("无法争取获取db信息，请检查配置")
		os.Exit(0)
	}
	return dbSlice
}

func checkDb(db string, engine *xorm.Engine)  {
	sql := fmt.Sprintf("use %s;", db)
	_, err := engine.QueryString(sql)
	if err != nil {
		// 说明访问数据库异常，可能不存在
		fmt.Println("数据库", db, "不存在,err is", err)
		os.Exit(0)
	}
}
