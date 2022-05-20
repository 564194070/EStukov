package db

import (
	"Web/db/mysql"
	"database/sql"
	"log"
)

func InsertFileMetaToDB(fileSha1 string, fileName string, fileSize int64, fileAddr string) bool {
	//IGNORE 忽略出错了行 将其他行插入到表中
	var err error
	defer func() {
		if p := recover(); p!= nil {
			log.Fatalln(p)
			log.Fatalln(err)

		}
	}()
	resInsert, err := mysql.DBConn().Prepare("insert IGNORE into file_information (`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values (?,?,?,?,1)")
	if err != nil {
		log.Fatalln("Failed to perpare statement, err", err)
		return false
	}

	defer resInsert.Close()

	resExec, err := resInsert.Exec(fileSha1, fileName, fileSize, fileAddr)
	if err != nil {
		log.Fatalln("Failed to exec Insert, err", err)
		return false
	}

	//看看受影响的行数
	rowLen, err := resExec.RowsAffected()
	if rowLen <= 0 {
		log.Fatalln("File Insert already hash = ",fileSha1)
		return false
	} else {
		return true
	}
	return false

}

type DBFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

//从Mysql获取文件源信息
func SelectFileMetaFromDB(fileSha1 string) (*DBFile, error)  {
	resSelect, err := mysql.DBConn().Prepare("select file_sha1,file_name,file_size,file_addr from file_information where file_sha1=? and status=1 limit 1")
	if err != nil {
		log.Fatalln("Failed to select from DB")
		return nil, err
	}
	defer resSelect.Close()

	resDBFile := new(DBFile)
	err = resSelect.QueryRow(fileSha1).Scan(&resDBFile.FileHash, &resDBFile.FileName, &resDBFile.FileSize, &resDBFile.FileAddr)
	if err != nil {
		log.Fatalln(err)
	}
	return resDBFile, nil
}