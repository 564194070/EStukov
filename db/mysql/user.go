package mysql

import (
	"log"
)

// UserInsertDB: 通过用户名以及密码完成user表的注册操作
func UserInsertDB(username string, passwd string) bool {
	resInsert, err := DBConn().Prepare("insert ignore into user_information (`user_name`,`user_pwd`) values(?,?)")
	if err != nil {
		log.Fatalln("Failed to insert User Information", err)
		return false
	}
	defer resInsert.Close()

	resExec, err := resInsert.Exec(username, passwd)
	if err != nil {
		log.Fatalln("Failed to Exec Insert User Information", resExec)
		return false
	}

	affected, err := resExec.RowsAffected()
	if err != nil && affected > 0 {
		return true
	}
	return false

}

func UserSelectSign(username string, passwd string) bool  {
	resSelect, err := DBConn().Prepare("Select user_pwd from user_information where user_name=? limit 1")
	if err != nil {
		log.Fatalln("Failed to select passwd from user table", err)
		return false
	}
	resQuery, err := resSelect.Query(username)
	if err != nil {
		log.Fatalln("Failed to query user info", err)
		return false
	} else if resQuery == nil{
		log.Fatalln("Failed to select user from database")
		return false
	}

	resRows := ParseRows(resQuery)
	if len(resRows) > 0 && string(resRows[0]["user_pwd"].([]byte)) == passwd {
		 return true
	}
	return false
}

// UpdateToken:刷新用户的Token
func UpdateToken(userName string, token string) bool {
	resUpdate, err := DBConn().Prepare("replace into user_token (`user_name`,`user_token`) values(?,?)")
	if err != nil {
		log.Fatalln("Failed to update token, err", err)
		return false
	}
	defer resUpdate.Close()

	_, err = resUpdate.Exec(userName,token)
	if err != nil {
		log.Fatalln("Failed to exec update, err", err)
		return false
	}
	return true
}

type User struct {
	UserName string
	Email string
	Phone string
	Signat string
	LastSign_at string
	Status int
}

func GetUserInfo(userName string) (User,error) {
	user := new(User)
	resSelect, err := DBConn().Prepare("select user_name,sign_at from user_information where user_name=?limit 1")
	if err != nil {
		log.Fatalln("Failed to select User Info, err", err)
		return *user, err
	}

	resSelect.QueryRow(userName).Scan(&user.UserName,&user.Signat)
	if err != nil {
		return *user, err
	}

	return *user, nil
}