package Ioc

import (
	"DentTick/Package/logger"
	"DentTick/Repository/Dao"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func InitDB(l logger.Logger) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}

	var cfg Config
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			//慢查询日志    0 为所有日志都打出来
			//慢查询阈值:200ms
			SlowThreshold: 200 * time.Millisecond,
			//日志等级:Error
			LogLevel: glogger.Error,
		}),
	})
	if err != nil {
		panic(err)
	}

	err = Dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

// 函数衍生类型实现接口
type gormLoggerFunc func(msg string, fields ...logger.Field)

func (f gormLoggerFunc) Printf(msg string, v ...interface{}) {
	f(msg, logger.Field{
		Key:   "args",
		Value: v,
	})
}
