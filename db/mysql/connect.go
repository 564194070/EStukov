package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" //将mysql的驱动注册到sql的context中被使用
	"log"
)

var database *sql.DB
var err error

func init()  {
	database, err = sql.Open("mysql", "EStukov:Mohican123.@tcp(119.3.232.83:3306)/File?charset=utf8")
	if err != nil {
		log.Fatalln("Mysql Failed to Open err",err)
	}

	//设置了最大连接数
	database.SetMaxIdleConns(10)
	err = database.Ping()
	if err != nil {
		log.Fatalln("Mysql Failed to connect err",err)
		panic(err)
	}
	fmt.Println("连接成功")
}

//返回数据库连接对象
func DBConn() *sql.DB {
	return database
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	//用户实现不知道返回了多少列 获取列表列名
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	// 迭代读取结果集合 是否还有未读取的数据记录
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		checkErr(err)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}