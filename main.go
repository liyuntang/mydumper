package main

import (
	"flag"
	"mydumper/common"
	"os"
)

// 定义参数变量
var (
	flagUser, flagPasswd, flagHost, flagOutDir, flagDbs, flagTabs, flagCharset string
	flagPort, flagThread int
	flagChunkSize, flagStmtSize int64
)

func init()  {
	flag.StringVar(&flagUser, "u", "root", "user")
	flag.StringVar(&flagPasswd, "p", "", "password")
	flag.StringVar(&flagHost, "h", "127.0.0.1", "db ip address")
	flag.StringVar(&flagOutDir, "o", "", "dump data dir")
	flag.StringVar(&flagDbs, "D", "", "dbs,db之间不能有空格，例如db1,db2,db3")
	flag.StringVar(&flagTabs, "T", "", "tables,table之间不能有空格，例如tab1,tab2,tab3")
	flag.StringVar(&flagCharset, "s", "utf8mb4", "charset")
	flag.IntVar(&flagPort, "P", 3306, "db port")
	flag.Int64Var(&flagChunkSize, "F", 10, "备份文件大小，单位M,默认128M")
	flag.Int64Var(&flagStmtSize, "S", 1, "每个insert大小，单位M,默认1M")
	flag.IntVar(&flagThread, "t", 4, "threads")
}

func main()  {
	flag.Parse()
	flagArgsLen := len(os.Args)
	if flagArgsLen <= 1 {
		flag.PrintDefaults()
	}

	// 将参数映射到struct
	flagArgs := &common.MyDumper{
		User:flagUser,
		Passwd:flagPasswd,
		Host:flagHost,
		Port:flagPort,
		Charset:flagCharset,
		Dbs:flagDbs,
		Tabs:flagTabs,
		OutDir:flagOutDir,
		ChunkSize:flagChunkSize,
		StmtSize:flagStmtSize,
		Threads:flagThread,
	}

	// 开始备份
	common.Start(flagArgs)
}
