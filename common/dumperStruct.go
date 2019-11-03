package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type MyDumper struct {
	User string
	Passwd string
	Host string
	Port int
	Charset string
	OutDir string
	Dbs string
	Tabs string
	ChunkSize int64
	StmtSize int64
	Threads int
}

func (md *MyDumper)getEngine() *xorm.Engine {
	endPoint := fmt.Sprintf("%s:%d", md.Host, md.Port)
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=true", md.User, md.Passwd, endPoint, "mysql", md.Charset)
	engine, err := xorm.NewEngine("mysql", dataSource)
	if err != nil {
		fmt.Println("init db connection", endPoint, "is bad")
		return nil
	}
	return engine
}

