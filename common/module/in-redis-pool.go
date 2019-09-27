package module

//需要进一步扩展
import (
	"fmt"
	"sync"

	"github.com/yamakiller/magicNet/library"
	"github.com/yamakiller/magicNet/util"
)

var (
	onceRedis sync.Once

	redisArray *redisDBArray
)

//redisDBArray Redis pool object
type redisDBArray struct {
	dbs map[int]*library.RedisDB
}

//RegisterDB Register db object
func (rda *redisDBArray) registerDB(host string, db int,
	maxIdle int,
	maxActive int,
	idleSec int) error {
	tmpdb := &library.RedisDB{}
	if err := tmpdb.Init(host, db, maxIdle, maxActive, idleSec); err != nil {
		return err
	}
	rda.dbs[db] = tmpdb
	return nil
}

//DB Retruns redis db object
func (rda *redisDBArray) db(db int) *library.RedisDB {
	r, success := rda.dbs[db]
	if !success {
		return nil
	}
	return r
}

func (rda *redisDBArray) close() {
	for _, v := range rda.dbs {
		v.Close()
	}
	rda.dbs = make(map[int]*library.RedisDB)
}

//redisInstance Redis connection pool interface
func redisInstance() *redisDBArray {
	onceRedis.Do(func() {
		redisArray = &redisDBArray{dbs: make(map[int]*library.RedisDB)}
	})

	return redisArray
}

//ReadisEnvAnalysis Environment variable parsing redis connector
func ReadisEnvAnalysis(m map[string]interface{}) error {
	redisEnv := util.GetEnvArray(m, "redis")
	if redisEnv == nil {
		return fmt.Errorf("Missing redis configuration information does not start properly")
	}

	for _, v := range redisEnv {
		one := util.ToEnvMap(v)
		host := util.GetEnvString(one, "host", "")
		db := util.GetEnvInt(one, "db", 0)
		maxIdle := util.GetEnvInt(one, "max-idle", 1)
		maxActive := util.GetEnvInt(one, "max-active", 1)
		idleSec := util.GetEnvInt(one, "idle-sec", 1000*60)

		if err := RedisRegister(host, db, maxIdle, maxActive, idleSec); err != nil {
			return err
		}
	}
	return nil
}

//RedisRegister Register Redis
func RedisRegister(host string, db int,
	maxIdle int,
	maxActive int,
	idleSec int) error {
	return redisInstance().registerDB(host, db, maxIdle, maxActive, idleSec)
}

//RedisDo Execute the Redis command
func RedisDo(db int, commandName string, args ...interface{}) (interface{}, error) {
	return redisInstance().db(db).Do(commandName)
}

//RedisClose Close the entire redis
func RedisClose() {
	redisInstance().close()
}
