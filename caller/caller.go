package caller

import (
	"fmt"

	"github.com/tydcwym/edgex_admin/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	// EdgexDB ...
	EdgexDB *gorm.DB
)

func InitClient() {
	// initRedisClient()
	initMysqlClient()
}

// func initRedisClient() {
// 	redisOpt := &redis.Options{
// 		Addr:     config.RedisConf.Address,
// 		Password: config.RedisConf.Password,
// 		DB:       config.RedisConf.DB,
// 	}
// 	RedisClient = redis.NewClient(redisOpt)
// }

func initMysqlClient() {

	optional := config.GetDefaultDBOptional()

	format := "%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=%s&readTimeout=%s&writeTimeout=%s"
	dbConfig := fmt.Sprintf(format, optional.User, optional.Password, optional.DBHostname, optional.DBPort,
		optional.DBName, optional.DBCharset, optional.Timeout, optional.ReadTimeout, optional.WriteTimeout)

	gormConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	var err error
	EdgexDB, err = gorm.Open(mysql.New(mysql.Config{
		DriverName: optional.DriverName,
		DSN:        dbConfig,
	}), &gormConfig)

	if err != nil {
		panic(err)
	}
}
