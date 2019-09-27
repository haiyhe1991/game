package module

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
	"testing"
)
func errCheck(err error) {
	if err != nil {
		fmt.Println("sorry,has some error:",err)
		os.Exit(-1)
	}
}
func Test_redisDBArray_registerDB(t *testing.T) {
	//使用redis封装的Dial进行tcp连接
	c,err := redis.Dial("tcp","localhost:6379")
	errCheck(err)

	defer c.Close()

	//对本次连接进行set操作
	_,setErr := c.Do("set","url","xxbandy.github.io")
	errCheck(setErr)

	//使用redis的string类型获取set的k/v信息
	r,getErr := redis.String(c.Do("get","url"))
	errCheck(getErr)
	fmt.Println(r)
}

func TestRedisDo(t *testing.T) {
//	redisInstance()
log.SetFlags(log.Llongfile|log.LstdFlags)
var err error
	if err = RedisRegister("127.0.0.1:6379",1,4,600000,600000);err != nil {
		log.Panic(err)
	}
	defer RedisClose()
	//TODO:如果注册为1号数据库  使用为0号数据
	log.Println(RedisDo(1,"set","url","xxbandy.github.io"))
}