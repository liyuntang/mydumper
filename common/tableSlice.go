package common

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"os"
	"strings"
)

var targe string = "liyuntangniubile"

func getTabSlice(dbSlice []string, tables string, engine *xorm.Engine) []string {
	dict := judgeTab(tables)
	// 根据dict的元素内容进行处理

	if strings.Trim(dict[0], "") == targe {
		fmt.Println("dump all")
		// 说明备份所有数据表，这时需要根据db获取table
		return dumpAll(dbSlice, engine)
	}
	// 说明需要备份一个或n多个表，这时候需要检测数据库与数据表是否对应
	tmpSlice := []string{}
	for _, db := range dbSlice {
		for _, tab := range dict {
			if checkTab(db, tab, engine) {
				// 说明数据表存在将其加入slice中
				//fmt.Println(db, tab, "is ok")
				tableInfo := fmt.Sprintf("%s.%s", db,tab)
				tmpSlice = append(tmpSlice, tableInfo)
			} else {
				fmt.Println(db, tab, "is bad")
				os.Exit(0)
			}
		}
	}
	return tmpSlice
}


func checkTab(db, table string, engine *xorm.Engine) bool {
	sql := fmt.Sprintf("desc %s.%s", db, table)
	_, err := engine.QueryString(sql)
	if err != nil {
		fmt.Println("db or table is not exist")
		return false
	}
	return true
}

func dumpAll(dbSlice []string, engine *xorm.Engine) []string {
	tableSlice := []string{}
	// 根据db信息获取该db下的所有数据表
	for _, db := range dbSlice {
		sql := fmt.Sprintf("select TABLE_NAME from information_schema.tables where TABLE_SCHEMA = '%s';", db)
		res, err := engine.QueryString(sql)
		if err != nil {
			fmt.Println("get table of db", db, "is bad")
			os.Exit(0)
		}
		for _, dict := range res {
			tableInfo := fmt.Sprintf("%s.%s", db, dict["TABLE_NAME"])
			tableSlice = append(tableSlice, tableInfo)
		}
	}
	return tableSlice
}

func judgeTab(tabs string) []string {
	// 声明一个slice用于存放table信息
	tabSlice := []string{}
	tmpSlice := []string{}
	for _, parame := range strings.Split(tabs, ",") {
		tmpSlice = append(tmpSlice, strings.Trim(parame, ""))
	}
	tmpLen := len(tmpSlice)
	if tmpLen == 1 {
		// 有两种可能：1、备份所有表      2、备份指定表
		if len(tmpSlice[0]) == 0 {
			// 备份所有表，如果是备份所有表的话，这里不显示表名，而是标记一下
			tabSlice = append(tabSlice, targe)
		} else {
			// 备份指定表,将这些表放入tabSlice中
			tabSlice = append(tabSlice, tmpSlice[0])
		}
	} else if tmpLen >= 2 {
		// 说明备份多个表
		for _, table := range tmpSlice {
			tabSlice = append(tabSlice, table)
		}
	} else {
		// 说明用户输入有误，提示报错，退出程序
		fmt.Println("-T参数输入有误，请核对参数")
		os.Exit(0)
	}
	return tabSlice
}
