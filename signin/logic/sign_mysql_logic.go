package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/yamakiller/magicNet/library"
	"log"
)

type UserInfo struct {
	Account  string
	Password string
	Origin   string
	UserId   int64
}

var MySql library.MySQLDB

const (
	USERNAME = "root"
	PASSWORD = "123456"
	NETWORK  = "tcp"
	SERVER   = "localhost"
	PORT     = 3306
	DATABASE = "game"
)

func init() {
	var (
		err error
	)
	MySql = library.MySQLDB{}
	dns := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)

	if err = MySql.Init(dns, 100, 50, 10); err != nil {
		log.Println(err)
		return
	}
}

func md5Data(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
func SetUser(account, password, origin string) (err error) {
	var ()
	//对密码进行md5
	password = md5Data(password)
	if _, err = MySql.Insert("insert INTO user(account,password,origin) values(?,?,?)", account, password, origin); err != nil {
		log.Println(err)
		return
	}
	return

}

func GetUserByAccount(account string) (info *UserInfo, err error) {
	var (
		result map[string]interface{}
	)
	info = &UserInfo{}
	if result, err = MySql.Query("SELECT id ,account,`password`,origin from  user WHERE account = ?", account); err != nil {
		log.Println(err)
		return
	}

	info.UserId = result["id"].(int64)
	info.Account = B2S(result["account"].([]uint8))
	info.Password = B2S(result["password"].([]uint8))
	info.Origin = B2S(result["origin"].([]uint8))

	return
}

func GetUserById(id int64) (info *UserInfo, err error) {
	var (
		result map[string]interface{}
	)
	info = &UserInfo{}
	if result, err = MySql.Query("SELECT id ,account,`password`,origin from  user WHERE id = ?", id); err != nil {
		log.Println(err)
		return
	}

	info.UserId = result["id"].(int64)
	info.Account = B2S(result["account"].([]uint8))
	info.Password = B2S(result["password"].([]uint8))
	info.Origin = B2S(result["origin"].([]uint8))

	return
}

func CheckUserByAccount(account string) (ret bool, err error) {
	var (
		result map[string]interface{}
	)
	if result, err = MySql.Query("SELECT  COUNT(*) FROM  user WHERE  account = ?", account); err != nil {
		log.Println(err)
		return
	}
	//log.Println(result["COUNT(*)"].(int64))
	if result["COUNT(*)"].(int64) == int64(1) {
		ret = true
	} else {
		ret = false
	}

	return
}

func CheckUserById(id int64) (ret bool, err error) {
	var (
		result map[string]interface{}
	)
	if result, err = MySql.Query("SELECT  COUNT(*) FROM  user WHERE  id = ?", id); err != nil {
		log.Println(err)
		return
	}
	//log.Println(result["COUNT(*)"].(int64))
	if result["COUNT(*)"].(int64) == int64(1) {
		ret = true
	} else {
		ret = false
	}

	return
}

func B2S(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}

//删除用户
func RemoveUser(account string, status int) (ret bool, err error) {
	var (
		//result map[string]interface{}
		n int64
	)
	if n, err = MySql.Update("UPDATE  `user` SET  `del` = ? WHERE account = ?", account, status); err != nil {
		log.Println(err)
		return
	}
	log.Println(n)

	return
}

func UpdateUserBanStauts(account string, status int) (ret bool, err error) {
	var (
		//result map[string]interface{}
		n int64
	)
	if n, err = MySql.Update("UPDATE  `user` SET  `ban` = ? WHERE account = ?", account, status); err != nil {
		log.Println(err)
		return
	}
	log.Println(n)
	return
}

//封账号
func UserBan(account string) (ret bool, err error) {
	return UpdateUserBanStauts(account, UserBanStatusOn)
}

//解封
func UserUnBan(account string) (ret bool, err error) {
	return UpdateUserBanStauts(account, UserBanStatusOff)
}
