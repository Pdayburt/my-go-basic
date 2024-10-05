//go:build !k8s

package config

var Config = config{
	MysqlConf: MysqlConfig{DSN: "root:root@tcp(localhost:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local"},
	RedisConf: RedisConfig{Addr: "localhost:6379"},
}
