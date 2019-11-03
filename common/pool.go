package common

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"time"
)

type conn struct {
	id int
	engine *xorm.Engine
}

type pool struct {
	conns chan *conn
}

func NewPool(flagArgs *MyDumper) *pool {
	// 由于pool是channel组成的，我们我们这里根据flagArgs.Threads的值来声明一个channel
	ch := make(chan *conn, flagArgs.Threads)
	// 使用MyDumper结构提的getEngine方法来获取数据库连接
	for i:=1;i<=flagArgs.Threads;i++ {
		Engine := flagArgs.getEngine()
		c :=  &conn{
			id:i,
			engine:Engine,
		}
		ch <- c
	}
	return &pool{
		conns:ch,
	}
}

func (p *pool)Get() *conn {
	for {
		select {
		case c := <- p.conns:
			fmt.Println("恭喜获取数据库连接成功,")
			return c
		case <- time.After(1*time.Second):
			fmt.Println("连接已经用完了，请耐心等待.............")
		}
	}

}

func (p *pool)Put(c *conn)  {
	p.conns <- c
}



















