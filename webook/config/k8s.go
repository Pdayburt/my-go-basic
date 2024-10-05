//go:build k8s

package config

var Config = config{
	MysqlConf: MysqlConfig{DSN: "root:root@tcp(webook-mysql:11309)/webook?charset=utf8mb4&parseTime=True&loc=Local"},
	RedisConf: RedisConfig{Addr: "webook-redis:11479"},
}
