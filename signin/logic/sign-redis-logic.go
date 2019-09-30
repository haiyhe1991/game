package logic

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/yamakiller/game/common/module"
	"log"
)

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}
func getUserSession(account string) {
}

const (
	UserStateLogin  = 1
	UserStateLogout = 0
)

type UserSession struct {
	State int    `json:"user_state"`
	Time   int64 `json:"log_time"`
}

var (
	ErrRedisMutexLockFailed   = errors.New("has been locked")
	ErrRedisMutexExpireFailed = errors.New("expire key failed")
	ErrRedisMutexUnlockFailed = errors.New("unlock failed")
)
var (
	ErrUserSetFaild = errors.New("user has  been exist")
	ErrUserSessionNotFind = errors.New("userSession not find")
)

func RedisMutexLock(account string) {
	var (
		err error
		ret interface{}
	)

	ret, err = module.RedisDo(1, "setnx", account+"lock", "lock")
	if err != nil {
		log.Println(err)
		return
	}
	//	log.Println(ret)
	if ret.(int64) == int64(1) {
		log.Println("success")
	} else {
		err = ErrRedisMutexLockFailed
		return
	}

	ret, err = module.RedisDo(1, "expire", account+"lock", 600)
	if err != nil {
		log.Println(err)
		return
	}
	if ret.(int64) == int64(1) {
		log.Println("success")
	} else {
		err = ErrRedisMutexExpireFailed
		return
	}
}

func RedisMutexUnlock(account string) {
	var (
		ret interface{}
	)
	ret, err := module.RedisDo(1, "del", account+"lock")
	if err != nil {
		log.Println(err)
		return
	}
	//	log.Println(ret)
	if ret.(int64) == int64(1) {
		log.Println("success")
	} else {
		err = ErrRedisMutexUnlockFailed
		return
	}

}

func SetUserSession(key string, session *UserSession) (err error) {
	var (
		b   []byte
		ret interface{}
	)
	if b, err = json.Marshal(session); err != nil {
		log.Println(err)
		return
	}
	ret, err = module.RedisDo(1, "SETNX", key, b)
	if err != nil {
		log.Println(err)
		return
	}
	if ret.(int64) == int64(1) {
		log.Println("success")
	} else {
		err = ErrUserSetFaild
		return
	}

	return
}

func GetUserSeesion(key string) (session *UserSession, err error) {
	var (
		//b []byte
		ret interface{}
	)
	session = &UserSession{}
	ret, err = module.RedisDo(1, "Get", key)
	if err != nil {
		log.Println(err)
		return
	}
	if ret == nil {
		err = ErrUserSessionNotFind
		return
	}

	if err = json.Unmarshal(ret.([]byte), session); err != nil {
		log.Println(err)
		return
	}

	return
}

func CheckUserSession(key string) (result bool) {

	var (
		ret interface{}
		err error
	)

	ret, err = module.RedisDo(1, "Get", key)
	if err != nil {
		log.Println(err)
		return
	}
	if ret == nil {
		result = false
		return
	}
	result = true
	return
}