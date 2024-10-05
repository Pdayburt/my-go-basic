package dao

import (
	"example.com/mod/webook/internal/domain"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestMySQL_Connect(t *testing.T) {
	InitViperV1()
	bd := InitDbB()
	var u domain.User
	bd.Find(&u, "id = ?", 1)
	fmt.Println(u)
}

// InitViperV1 使用结构体读取额配置文件
func InitViperV1() {
	viper.SetConfigFile("/Users/anatkh/Downloads/blockChain/golang/my-go-basic/webook/config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
func InitDbB() *gorm.DB {

	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config

	err := viper.UnmarshalKey("db.mysql", &cfg)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	/*err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}*/
	return db
}
