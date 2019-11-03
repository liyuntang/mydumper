package common

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

func dumpTabData(outDir, tableInfo string, chunkSize, stmtSize int64, engine *xorm.Engine, rows int)  {
	//fmt.Println("开始备份表数据, tableInfo is", tableInfo)
	// 将 chunkSize、stmtSize 转换成byte单位
	fileSize := chunkSize * 1024 * 1024
	insertSize := stmtSize * 1024 *1024
	num, err := strconv.Atoi(strconv.FormatInt(insertSize, 10))
	if err != nil {
		fmt.Println("get dumpFile size is bad, err is", err)
		os.Exit(0)
	}
	// 获取表结构信息，并生成有序的slice
	clumonSlice := getClumon(engine, tableInfo)

	//查询数据
	startID := 1
	endID := 10000
	step :=10000
	for startID <= rows {
		sql := fmt.Sprintf("SELECT /*backup*/ * FROM %s where id between %d and %d;", tableInfo, startID, endID)
		resultSlice, err := engine.QueryString(sql)
		if err != nil {
			fmt.Println("get data from table is bad, table is", tableInfo, "err is", err)
			os.Exit(0)
		}
		sql_1 := strings.Join(clumonSlice, ",")
		sql_2 := fmt.Sprintf("insert into %s(%s", tableInfo, sql_1)
		sql_3 := fmt.Sprintf("%s) values ", sql_2)
		for _, dict := range resultSlice {
			data := insertSql(clumonSlice, dict)
			//fmt.Println("len of sql_3 is", len(sql_3), "num is", num)
			if len(sql_3) >= num {
				sql_3 += strings.Trim(data, ",")+";\r\n"
				dumpFile := fileName(outDir, tableInfo, fileSize)
				writeData(dumpFile, sql_3)
				sql_3 = fmt.Sprintf("%s) values ", sql_2)
			} else {
				sql_3 += data
			}

		}
		startID += step
		endID += step
	}
	fmt.Println("dump table is ok, table is", tableInfo)
}

func insertSql(clumonSlice []string, data map[string]string) string {
	var value string
	num := len(clumonSlice)-1
	for index, clumon := range clumonSlice {
		if index == num {
			value += data[clumon]
		} else {
			value += data[clumon]+","
		}

	}
	sql_4 := fmt.Sprintf("（%s),", value)
	return sql_4
}

func fileName(outDir, tableInfo string, fileSize int64) string {
	// 转换名称
	tmpDT := strings.Replace(tableInfo, ".", "-",1)
	// 备份文件名称
	//var dumpFile string
	// 声明一个存放备份文件名称的slice，改slice是有序的
	fileSlice := []int{}
	dirInfo, err := ioutil.ReadDir(outDir)
	if err != nil {
		fmt.Println("list outDir", outDir, "is bad, err is", err)
		os.Exit(0)
	}
	for _, dir := range dirInfo {
		f := dir.Name()
		//fmt.Println("tableInfo is", tableInfo)
		// 判断f是不是备份文件，依据有三个：1、不能以-create-schema.sql结尾      2、不能以-create-table.sql结尾  3、包含tmpDT
		if !strings.HasSuffix(f, "-create-schema.sql") && !strings.HasSuffix(f,"-create-table.sql") && strings.Contains(f, tmpDT) {
			// 说明不是备份文件
			// 说明该文件是备份文件，此时备份文件可能有多个也可能只有一个，需要兼容

			n, err := strconv.Atoi(strings.Split(strings.Split(f, "-")[2], ".")[0])
			if err != nil {
				fmt.Println("err is", err)
				os.Exit(0)
			}
			fileSlice = append(fileSlice, n)
		}
	}
	// 判断fileSlice是否有元素，如果有则说明存在备份文件，如果为空则表明没有备份文件，此时给个初始化的备份文件即可
	Len := len(fileSlice)
	if Len == 0 {
		// 说明当前不存在备份数据文件，此时给个初始化的备份文件即可
		return fmt.Sprintf("%s/%s-1.sql", outDir, tmpDT)
	}
	// 说明当前存在备份文件，此时需要将fileSlice变成有序的
	sort.Ints(fileSlice)
	// 根据fileSlice里的文件编号取出当前文件
	m := fileSlice[Len-1]
	nowDumpFile := fmt.Sprintf("%s/%s-%d.sql", outDir, tmpDT, m)
	// 此时需要判断备份文件的大小，如果>= fileSize则创建新文件,新文件名称要累加，否则返回当前文件名称
	fileInfo, err1 := os.Stat(nowDumpFile)
	if err1 != nil {
		fmt.Println("获取nowDumpFile状态失败,err is", err1)
		os.Exit(0)
	}
	if fileInfo.Size() >= fileSize {
		// 需要切割文件
		return fmt.Sprintf("%s/%s-%d.sql", outDir, tmpDT, m+1)
	} else {
		// 不需要切割文件，返回f即可
		return nowDumpFile
	}
}

func getClumon(engine *xorm.Engine, tableInfo string) []string {
	// 定义一个slice用于存放表的字段
	clumonSlice := []string{}
	//获取表结构
	sql := fmt.Sprintf("DESC %s;", tableInfo)
	resultSlice, _ := engine.QueryString(sql)
	for _, dict := range resultSlice {
		clumonSlice = append(clumonSlice, dict["Field"])
	}
	sort.Strings(clumonSlice)
	return clumonSlice
}


func getRows(engine *xorm.Engine, tableInfo string) int {
	sql := fmt.Sprintf("select count(id) from %s;", tableInfo)
	res, _ := engine.QueryString(sql)
	num, _ := strconv.Atoi(res[0]["count(id)"])
	return num
}







