package fileupload

import (
	"Web/cache/redis"
	"Web/db"
	"Web/db/mysql"
	"Web/util"
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

//MultipartUploadInfo 初始化信息  
type MultipartUploadInfo struct {
	FileHash string
	FileSize int
	UploadID string
	ChunkSize int
	ChunkCount int
}

//初始化分块上传
func InitMultipartUploadHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	userName := r.Form.Get("userName" )
	fileHash := r.Form.Get("fileHash" )
	fileSize, err:= strconv.Atoi(r.Form.Get("fileSize" ))
	if err != nil {
		w.Write(util.NewReqMessage(-1,"Failed use",nil).JsonToByte())
	}

	//获得redis连接
	rConn := redis.RedisPool()
	defer rConn.Close()
	//生成分块上传的初始化信息
	upInfo := MultipartUploadInfo{
		FileHash: fileHash,
		FileSize: fileSize,
		UploadID: userName + fmt.Sprintf("%x",time.Now().UnixNano()),
		ChunkSize: 5*1024*1024,
		ChunkCount: int(math.Ceil(float64(fileSize)/(5*1024*1024))),
	}

	//生成唯一上传ID
	//将初始化信息写入redis
	rConn.HSet("MP_"+upInfo.UploadID,"chunkcount",upInfo.ChunkCount)
	rConn.HSet("MP_"+upInfo.UploadID,"filehash",upInfo.FileHash)
	rConn.HSet("MP_"+upInfo.UploadID,"filesize",upInfo.FileSize)

	//缓存分块初始化信息
	//将响应信息返回给客户端
	w.Write(util.NewReqMessage(0,"OK",nil).JsonToByte() )
}

// UploadPartHandler : 上传文件分块
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求参数
	r.ParseForm()
	//	username := r.Form.Get("username")
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	// 2. 获得redis连接池中的一个连接
	rConn := redis.RedisPool()
	defer rConn.Close()

	// 3. 获得文件句柄，用于存储分块内容
	fpath := "/data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		w.Write(util.NewReqMessage(-1, "Upload part failed", nil).JsonToByte())
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 4. 更新redis缓存状态
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 5. 返回处理结果到客户端
	w.Write(util.NewReqMessage(0, "OK", nil).JsonToByte())
}

//CompleteUploadHandler:通知上传合并接口
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	upid := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("fileName")

	//获得redis连接
	rConn := redis.RedisPool()
	defer rConn.Close()
	//通过uoloadid哦按段所有分块是否上传完成
	result, err := rConn.HGetAll("MP" + upid).Result()
	if err != nil {
		w.Write(util.NewReqMessage(-1,"failed to upload all chunk",nil).JsonToByte())
		return
	}
	totalCount, chunkCount := 0, 0
	for field, val := range result {
		if field == "chunkcount" {
			totalCount, _ = strconv.Atoi(val)
		} else if strings.HasPrefix(field, "chkidx_") && val == "1" {
			chunkCount ++
		}
	}
	if totalCount != chunkCount {
		w.Write(util.NewReqMessage(-2,"invalid request", nil).JsonToByte())
		return
	}


	//合并分块
	//更新唯一文件表
	fsize, _ := strconv.Atoi(filesize)
	db.InsertFileMetaToDB(filehash,filename,int64(fsize),"")
	//更新用户文件表
	mysql.InsertFileToUserFile(username,filehash,filename,int64(fsize))
	//响应处理结果
}

