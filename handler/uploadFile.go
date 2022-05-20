package fileupload

import (
	"Web/db"
	"Web/db/mysql"
	"Web/meta"
	"Web/util"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)


// FileUploadHandler 处理单一文件上传和返回界面的逻辑
func FileUploadHandler(w http.ResponseWriter,r *http.Request) {
	if r.Method == http.MethodGet {
		//返回用户上传的文件
		data, err := ioutil.ReadFile("./static/upload.html")
		if err != nil {
			io.WriteString(w, "ReadFile False")
		}
		io.WriteString(w, string(data))
		//c.Writer.WriteString()
	} else if r.Method == "POST" {
		//接受用户文件流存储到本地目录
		//文件句柄 文件头 错误信息
		file, fileHead, err := r.FormFile("uploadfile")
		if err != nil {
			log.Fatalln("Failed to get data, err: ", err)
			return
		}
		defer file.Close()


		fileMeta := meta.FileMeta{
			FileName:     fileHead.Filename,
			LocationPath: "./tmp/" + fileHead.Filename,
			UploadTime:   time.Now().Format("2016-01-02 15:04:05"),
		}

		//写的路径 + 文件名
		filePoint, err := os.Create(fileMeta.LocationPath)
		if err != nil {
			log.Fatalln("Failed to create file, err: ", err)
			return
		}
		defer filePoint.Close()
		//将用户文件Copy到本地到新文件中
		fileMeta.FileSize, err = io.Copy(filePoint, file)
		if err != nil {
			log.Fatalln("Failed to save data into file, err:", err)
		}

		filePoint.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(filePoint)
		//meta.UpdataFileMetaToMap(fileMeta)
		//写唯一文件表
		meta.UpdateFileMetaToDB(fileMeta)

		//写用户文件表
		r.ParseForm()
		userName := r.Form.Get("userName")
		isInsert := mysql.InsertFileToUserFile(userName, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
		if isInsert != true {
			log.Fatalln("Failed to insert, err")
			w.Write([]byte("error to insert"))
		} else {
			http.Redirect(w,r, "/file/upload/suc",http.StatusFound)
		}

		//http.Redirect(c.Writer, c.Request, "/file", http.StatusFound)
	}
}

// UploadSuccessHandler 成功之后重定向到此处
func UploadSuccessHandler(w http.ResponseWriter,r *http.Request) {
	io.WriteString(w,"Ok")
}

//GetFileMetaHandler：获取文件的源信息
func GetFileMetaHandler(w http.ResponseWriter,r *http.Request) {
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	//userFileMeta := meta.GetFileMetaFromMap(filehash)
	userFileMeta, err := db.SelectFileMetaFromDB(filehash)
	if err != nil {
		log.Fatalln("Failed to get Meta ", err)
	}
	//从结构体转换成json string的格式返回到客户端
	userFileMetaData, err := json.Marshal(userFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(userFileMetaData)
}

//DownloadHandler：文件下载接口
func DownloadHandler(w http.ResponseWriter,r *http.Request)  {
	r.ParseForm()

	filehash := r.Form.Get("filehash")
	//加载已存储到本地云端的文件内容，并返回到客户端 就像静态资源

	//userFileMeta := meta.GetFileMetaFromMap(filehash)
	userFileMeta, err := db.SelectFileMetaFromDB(filehash)
	if err != nil {
		log.Fatalln("Failed to get Meta From base", err)
	}
	filePoint, err := os.Open(userFileMeta.FileAddr.String)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer filePoint.Close()

	//如果文件大 就流状读取 读一部分 刷一部分
	userFileData, err := ioutil.ReadAll(filePoint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type","application/octect-stream")
	w.Header().Set("Content-Descrption","attachment;filename=\""+userFileMeta.FileName.String+"\"")
	w.Write(userFileData)
}

// FileUpdateMetaHandler:修改文件名
func FileUpdateMetaHandler(w http.ResponseWriter,r *http.Request) {
	r.ParseForm()

	behavior := r.Form.Get("behavior")
	userFileHash := r.Form.Get("filehash")
	newFileName := r.Form.Get("fileName")

	switch behavior {
	case "0 ":
		userFileMeta := meta.GetFileMetaFromMap(userFileHash)
		userFileMeta.FileName = newFileName
		meta.UpdataFileMetaToMap(userFileMeta)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userFileData, err := json.Marshal(userFileHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userFileData)
}


// FileDelHandler:删除文件
func FileDelHandler(w http.ResponseWriter,r *http.Request)  {
	r.ParseForm()

	userFileHash := r.Form.Get("filehash")
	meta.RemoveFileMeta(userFileHash)

	userFileMeta := meta.GetFileMetaFromMap(userFileHash)
	os.Remove(userFileMeta.LocationPath)
	w.WriteHeader(http.StatusOK)
}

//批量查询文件源信息
//GetFilesMetaHandler：获取文件的源信息
func GetFilesMetaHandler(w http.ResponseWriter,r *http.Request) {
	r.ParseForm()
	limitFile, _ := strconv.Atoi(r.Form.Get("limit"))
	userName := r.Form.Get("userName")
	resSelect, err := mysql.SelectUserFileMetas(userName, limitFile)

	//filehash := r.Form["filehash"][0]
	//userFileMeta := meta.GetFileMetaFromMap(filehash)
	//userFileMeta, err := db.GetFileMetaFromDB(filehash)

	if err != nil {
		log.Fatalln("Failed to get Meta ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//从结构体转换成json string的格式返回到客户端
	userFileMetaData, err := json.Marshal(resSelect)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(userFileMetaData)
}

// TryFastUploadHandler:尝试秒传接口
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	//解析请求参数
	userName := r.Form.Get("userName")
	fileHash := r.Form.Get("fileHash")
	fileName := r.Form.Get("fileName")
	fileSize, _ := strconv.Atoi(r.Form.Get("fileSize"))
	//从文件表中查询相同hash的文件记录)
	fileMeta, err := db.SelectFileMetaFromDB(fileHash)
	if err != nil {
		log.Fatalln("Failed to select meta, err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//查不到记录返回秒传失败
	if fileMeta == nil {
		resp := util.JsonMessage{
			ResCode: -1,
			Message: "Failed to fast upload, restart upload normal",
		}
		w.Write(resp.JsonToByte())
		return
	}
	//上传过就直接将信息写入用户文件表 返回成功
	resInset := mysql.InsertFileToUserFile(userName, fileHash, fileName, int64(fileSize))
	if true == resInset {
		resp := util.JsonMessage{
			ResCode: 0,
			Message: "Fast upload ok",
		}
		w.Write(resp.JsonToByte())
		return
	} else {
		resp := util.JsonMessage{
			ResCode: -2,
			Message: "Fast upload failed",
		}
		w.Write(resp.JsonToByte())
		return
	}
}