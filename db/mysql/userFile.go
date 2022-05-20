package mysql

import (
	"log"
	"time"
)

// UserFile:用户文件表结构体
type UserFile struct {
	UserName string
	FileHash string
	FileName string
	FileSize int64
	UploadAt string
	LastUpdated string
}

// InsertFileToUserFile:插入用户文件表
func InsertFileToUserFile(userName, fileHash, fileName string, fileSize int64) bool {
	resInsert, err := DBConn().Prepare("insert ignore into user_file_information(`user_name`,`user_file_sha1`,`user_file_size`,`upload_at`) values (?,?,?,?,?)")
	if err != nil {
		log.Fatalln("Failed to insert in userFile table, err", err)
	}
	defer resInsert.Close()
	_, err = resInsert.Exec(userName, fileHash, fileName, fileSize, time.Now())
	if err != nil {
		log.Fatalln("Failed to exec insert userFile table, err", err)
		return false
	}
	return true
}


func SelectUserFileMetas(userName string, limit int) ([]UserFile, error) {
	resSelect, err := DBConn().Prepare("select user_file_sha1,user_file_size,upload_at,last_update from user_file_information where user_name=?limit ?")
	if err != nil {
		log.Fatalln("Failed to select, err", err)
		return nil, err
	}

	resRow, err := resSelect.Query(userName, limit)
	if err != nil {
		log.Fatalln("Fialed to Query, err", err)
		return nil, err
	}

	var userFiles []UserFile
	for resRow.Next() {
		userfile := new(UserFile)
		err := resRow.Scan(userfile.FileHash, userfile.FileName, userfile.FileSize, userfile.UploadAt, userfile.LastUpdated)
		if err != nil {
			log.Fatalln("Failed to Scan")
			break
		}
		userFiles = append(userFiles, *userfile)
	}
	return userFiles, nil
}