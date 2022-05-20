package meta

import (
	"Web/db"
	"log"
)

//FileMeta: 文件源信息结构
type FileMeta struct {
	FileSha1     string //文件的唯一标志
	FileName     string
	FileSize     int64
	LocationPath string
	UploadTime   string
}

//通过文件来查找源信息
var sha1FileSearchMap map[string]FileMeta

func init()  {
	sha1FileSearchMap = make(map[string]FileMeta, 10)
}
//UpdataFileMetaToMap: 新增/更新文件源信息
func UpdataFileMetaToMap(fileMeta FileMeta) {
	sha1FileSearchMap[fileMeta.FileSha1] = fileMeta

}
//UpdateFileMetaToDB: 新增/更新文件源信息到Mysql中
func UpdateFileMetaToDB(fileMeta FileMeta)  {
	resInsert := db.InsertFileMetaToDB(fileMeta.FileSha1,fileMeta.FileName,fileMeta.FileSize,fileMeta.LocationPath)
	if resInsert != true {
		log.Fatalln("Insert Error",resInsert)
	}
}

//GetFileMeta: 通过SHA1获取文件源信息对象
func GetFileMetaFromMap(fileSha1 string) FileMeta {
	return sha1FileSearchMap[fileSha1]
}
// GetFileMetaFromDB: 通过数据库获取文件的源信息
func GetFileMetaFromDB(fileSha1 string) (FileMeta, error) {
	resFileMeta, err := db.SelectFileMetaFromDB(fileSha1)
	if err != nil {
		return FileMeta{}, nil
	}
	fileMeta := FileMeta{
		FileSha1: resFileMeta.FileHash,
		FileSize: resFileMeta.FileSize.Int64,
		FileName: resFileMeta.FileName.String,
		LocationPath: resFileMeta.FileAddr.String,
	}
	return fileMeta, nil
}
//删除源信息
func RemoveFileMeta(fileSha1 string)  {
	delete(sha1FileSearchMap,fileSha1)
}