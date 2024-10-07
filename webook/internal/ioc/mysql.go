package ioc

import (
	"example.com/mod/webook/internal/repository/dao"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDbB() *gorm.DB {
	//InitViperV1()
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config

	err := viper.UnmarshalKey("db.mysql", &cfg)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func InitViperV1() {
	viper.SetConfigFile("/Users/anatkh/Downloads/blockChain/golang/my-go-basic/webook/config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
