package config

type config struct {
	MysqlConf MysqlConfig
	RedisConf RedisConfig
}

type MysqlConfig struct {
	DSN string
}
type RedisConfig struct {
	Addr string
}
