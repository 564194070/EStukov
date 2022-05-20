package fileupload

import (
	"Web/db/mysql"
	"Web/util"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const saltOfPasswd string = "Mohican"

// UserRegisterHandler:处理用户注册请求
func UserRegisterHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/register.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	} else if  r.Method == http.MethodPost {
		r.ParseForm()
		userName := r.Form.Get("userName")
		userPwd := r.Form.Get("passWd")
		fmt.Println(r.Form)
		fmt.Println("开始了")

		if len(userName) < 5 || len(userPwd) < 5 {
			luN,luP := len(userName), len(userPwd)
			str := "string of userName or passwd should be more" + strconv.Itoa(luN)+strconv.Itoa(luP)
			w.Write([]byte(str))
			return
		}
		passwd := util.Sha1([]byte(userPwd + saltOfPasswd))
		resInsert := mysql.UserInsertDB(userName,passwd)
		if true == resInsert {
			w.Write([]byte("SUCCESS"))
		} else {
			w.Write([]byte("FAILED"))
		}

	}

	type Task interface {
		Init()
		Run()
	}
	type manager map[string]Task

	var mgr = new(manager)
	(*mgr)["name"] = "10"

}

// UserLoginHandler:用户登录函数
func UserLoginHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/login.html")
		if err != nil {
			 	w.WriteHeader(http.StatusInternalServerError)
				return
			}
		w.Write(data)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		userName := r.Form.Get("username")
		userPwd := r.Form.Get("password")
		passWord := util.Sha1([]byte(userPwd + saltOfPasswd))

		//认证
		userChecked := mysql.UserSelectSign(userName,passWord)
		if userChecked == false {
			w.Write([]byte("Failed to login"))
		}

		token := MakeToken(userName)
		updateToken := mysql.UpdateToken(userName, token)
		if updateToken  == false {
			w.Write([]byte("Failed to update token"))
		}

		//鉴权
		//重定向 返回用户Token
		resp := util.JsonMessage{
			ResCode: 0,
			Message: "Login Successful",
			Extend: struct {
				Location string
				UserName string
				Token string
			}{
				Location: "http://" + r.Host + "/user/index",
				UserName: userName,
				Token: token,
			},
		}
		w.Write(resp.JsonToByte())
	}

}

// UserInfoHandler:查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request)  {
	//解析请求参数
	r.ParseForm()
	userName := r.Form.Get("username")
	token := r.Form.Get("token")
	//验证token是否有效
	tokenUseful := IsTokenVaild(token)
	if tokenUseful == false {
		w.WriteHeader(http.StatusForbidden)
	}
	//查询用户信息
	userInfo, err := mysql.GetUserInfo(userName)
	if err != nil {
		log.Fatalln("Failed to get user info , err", err)
		w.WriteHeader(http.StatusForbidden)
	}
	//响应
	rep := util.JsonMessage{
		ResCode: 0,
		Message: "OK",
		Extend: userInfo,
	}
	w.Write(rep.JsonToByte())
}
func MakeToken(userName string) string {
	//md5(userName + timestamp + tokensalt) + timestamp[:8]
	timestamp := fmt.Sprint("%x",time.Now().Unix())
	tokenPrefix := util.MD5([]byte(userName + timestamp + "tokenMohican"))
	return tokenPrefix + timestamp[:8]
}

func IsTokenVaild (token string) bool {
	if 40 != len(token) {
		return false
	}
	//时间

	return true
}

func UserIndexHandler(w http.ResponseWriter, r *http.Request)  {
	data, err := ioutil.ReadFile("./static/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}