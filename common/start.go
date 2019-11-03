package common

import (
	"fmt"
	"sync"
)
var (
	wg sync.WaitGroup
)

func Start(flagArgs *MyDumper)  {
	// 判断 flagArgs.OutDir的状态
	judgeOutDir(flagArgs.OutDir)

	// 生成连接池
	pool := NewPool(flagArgs)

	// 根据 flagArgs.Dbs、flagArgs.Tabs的设置生成dbSlice及tabSlice，用于后面的备份
	dbconn := pool.Get()
	dbSlice := getDbSlice(flagArgs.Dbs, dbconn.engine)
	tabSlice := getTabSlice(dbSlice, flagArgs.Tabs, dbconn.engine)
	//fmt.Println("tableSlice is", tabSlice)
	// 单线程备份schema信息,备份完成后将连接返回连接池
	dumpDb(dbSlice, flagArgs.OutDir, dbconn.engine)
	pool.Put(dbconn)

	// 多线程备份table信息、表数据
	for i:=0;i<len(tabSlice);i++{
		wg.Add(1)
		dbConn := pool.Get()
		fmt.Println("conn id is", dbConn.id)
		tableInfo := tabSlice[i]
		go func() {
			// 吃水不忘挖井人
			defer pool.Put(dbConn)
			// 备份表结构信息
			dumpTabInfo(flagArgs.OutDir, tableInfo, dbConn.engine)
			// 备份表数据
			// 查询数据表是否是空表，如果是空表则不用备份，如果不为空则需要备份

			rows := getRows(dbConn.engine, tableInfo)
			if rows == 0{
				// 说明是空表，不用备份，直接返回即可
				wg.Done()
			}
			// 说明数据表不为空，需要备份表数据
			dumpTabData(flagArgs.OutDir, tableInfo, flagArgs.ChunkSize, flagArgs.StmtSize, dbConn.engine, rows)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("over.......................")
}

